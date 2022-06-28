package routes

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/foolin/pagser"

	"path/filepath"

	"scrapper/src/mongodb"
	"scrapper/src/types"
	"scrapper/src/utils"

	"github.com/jinzhu/copier"

	"golang.org/x/exp/slices"

	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"regexp"
	"strings"
)

type ParamsSearchPhotoFlickr struct {
	Quality string `uri:"quality" binding:"required"`
}

// Find all the photos with specific quality and folder directory.
func SearchPhotosFlickr(mongoClient *mongo.Client, params ParamsSearchPhotoFlickr) ([]primitive.ObjectID, error) {

	quality := params.Quality
	var insertedIDs []primitive.ObjectID

	parser := pagser.New() // parsing html in string responses

	// If path is already a directory, MkdirAll does nothing and returns nil
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	origin := "flickr"
	err := os.MkdirAll(filepath.Join(folderDir, origin), os.ModePerm)
	if err != nil {
		return nil, err
	}

	collectionImages:= mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))

	unwantedTags, err := mongodb.TagsUnwantedNames(mongoClient)
	if err != nil {
		return nil, err
	}
	sort.Strings(unwantedTags)

	wantedTags, err := mongodb.TagsWantedNames(mongoClient)
	if err != nil {
		return nil, err
	}
	sort.Strings(wantedTags)

	for _, wantedTag := range wantedTags {

		// all the commercial use licenses
		// https://www.flickr.com/services/api/flickr.photos.licenses.getInfo.html
		var licenseIDsNames = map[string]string{
			"4":  "Attribution License",
			"5":  "Attribution-ShareAlike License",
			"7":  "No known copyright restrictions",
			"9":  "Public Domain Dedication (CC0)",
			"10": "Public Domain Mark",
		}
		licenseIDs := [5]string{"4", "5", "7", "9", "10"}
		for _, licenseID := range licenseIDs {

			// start with the first page
			page := 1
			searchPerPage, err := searchPhotosPerPageFlickr(parser, licenseID, wantedTag, strconv.FormatUint(uint64(page), 10))
			if err != nil {
				return nil, fmt.Errorf("searchPhotosPerPageFlickr has failed: \n%v", err)
			}

			for page := page; page <= int(searchPerPage.Pages); page++ {
				searchPerPage, err := searchPhotosPerPageFlickr(parser, licenseID, wantedTag, strconv.FormatUint(uint64(page), 10))
				if err != nil {
					return nil, fmt.Errorf("searchPhotosPerPageFlickr has failed: \n%v", err)
				}
				for _, photo := range searchPerPage.Photos {

					// look for existing image
					_, err := mongodb.FindImageIDByOriginID(collectionImages, photo.ID)
					if err != nil {
						return nil, err
					}

					// extract the photo informations
					infoData, err := infoPhoto(parser, photo)
					if err != nil {
						return nil, fmt.Errorf("InfoPhoto has failed: %v", err)
					}
					var photoTags []string
					for _, tag := range infoData.Tags {
						photoTags = append(photoTags, strings.ToLower(tag.Name))
					}

					// skip image if one of its tag is unwanted
					idx := utils.FindIndexRegExp(unwantedTags, photoTags)
					if idx != -1 {
						continue // skip image with unwated tag
					}

					// extract the photo download link
					downloadData, err := downloadPhoto(parser, photo.ID)
					if err != nil {
						return nil, fmt.Errorf("DownloadPhoto has failed: \n%v", err)
					}

					// get the download link for the correct resolution
					label := strings.ToLower(quality)
					regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, label)
					idx = slices.IndexFunc(downloadData.Photos, func(download DownloadPhotoSingleData) bool { return strings.ToLower(download.Label) == label })
					if idx == -1 {
						idx = slices.IndexFunc(downloadData.Photos, func(download DownloadPhotoSingleData) bool {
							matched, err := regexp.Match(regexpMatch, []byte(strings.ToLower(download.Label)))
							if err != nil {
								return false
							}
							return matched
						})
					}
					if idx == -1 {
						return nil, fmt.Errorf("Cannot find label %s and its derivatives %s in SearchPhoto! id %s has available the following:\n%v\n", label, regexpMatch, photo.ID, downloadData)
					}

					// download photo into folder and rename it <id>.<format>
					fileName := fmt.Sprintf("%s.%s", photo.ID, infoData.OriginalFormat)
					path := fmt.Sprintf(filepath.Join(folderDir, origin, fileName))
					err = DownloadFile(downloadData.Photos[idx].Source, path)
					if err != nil {
						return nil, err
					}

					// tags creation
					var tags []types.Tag
					copier.Copy(&tags, &infoData.Tags)
					now := time.Now()
					for i := 0; i < len(tags); i++ {
						tag := &tags[i]
						tag.Name = strings.ToLower(tag.Name)
						tag.CreationDate = &now
						tag.Origin = origin
					}

					// image creation
					document := types.Image{
						Origin:       origin,
						OriginID:     photo.ID,
						Extension:    infoData.OriginalFormat,
						Path:         fileName,
						Width:        downloadData.Photos[idx].Width,
						Height:       downloadData.Photos[idx].Height,
						Title:        infoData.Title,
						Description:  infoData.Description,
						License:      licenseIDsNames[licenseID],
						CreationDate: &now,
						Tags:         tags,
					}

					insertedID, err := mongodb.InsertImage(collectionImages, document)
					if err != nil {
						return nil, err
					}
					insertedIDs = append(insertedIDs, insertedID)
				}
			}
		}
	}
	return insertedIDs, nil
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type SearchPhotPerPageData struct {
	Stat    string  `pagser:"rsp->attr(stat)"`
	Page    uint    `pagser:"photos->attr(page)"`
	Pages   uint    `pagser:"photos->attr(pages)"`
	PerPage uint    `pagser:"photos->attr(perpage)"`
	Total   uint    `pagser:"photos->attr(total)"`
	Photos  []PhotoFlickr `pagser:"photo"`
}
type PhotoFlickr struct {
	ID     string `pagser:"->attr(id)"`
	Secret string `pagser:"->attr(secret)"`
	Title  string `pagser:"->attr(title)"`
}

// Search images for one page of max 500 images
func searchPhotosPerPageFlickr(parser *pagser.Pagser, ids string, tags string, page string) (*SearchPhotPerPageData, error) {
	r := &Request{
		Host: "https://api.flickr.com/services/rest/?",
		Args: map[string]string{
			"api_key": utils.DotEnvVariable("FLICKR_PRIVATE_KEY"),
			"method":  "flickr.photos.search",
			"tags":    tags,
			"license": ids,
			"media":   "photos",
			"page":    page,
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, err
	}

	var pageData SearchPhotPerPageData
	err = parser.Parse(&pageData, string(body))
	if err != nil {
		return nil, err
	}
	fmt.Println(utils.ToJSON(pageData))
	if pageData.Stat != "ok" {
		return nil, fmt.Errorf("SearchPhotoPerPageRequest is not ok\n%v\n", pageData)
	}
	if pageData.Page == 0 || pageData.Pages == 0 || pageData.PerPage == 0 || pageData.Total == 0 {
		return nil, errors.New("Some informations are missing from SearchPhotoPerPage")
	}
	return &pageData, nil
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type DownloadPhotoSingleData struct {
	Label  string `pagser:"->attr(label)"`
	Width  int   `pagser:"->attr(width)"`
	Height int   `pagser:"->attr(height)"`
	Source string `pagser:"->attr(source)"`
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type DownloadPhotoData struct {
	Stat   string                    `pagser:"rsp->attr(stat)"`
	Photos []DownloadPhotoSingleData `pagser:"size"`
}

func downloadPhoto(parser *pagser.Pagser, id string) (*DownloadPhotoData, error) {
	r := &Request{
		Host: "https://api.flickr.com/services/rest/?",
		Args: map[string]string{
			"api_key":  utils.DotEnvVariable("FLICKR_PRIVATE_KEY"),
			"method":   "flickr.photos.getSizes",
			"photo_id": id,
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, fmt.Errorf("DownloadPhoto has failed: \n%v", err)
	}

	var downloadData DownloadPhotoData
	err = parser.Parse(&downloadData, string(body))
	if err != nil {
		return nil, err
	}
	fmt.Println(utils.ToJSON(downloadData))

	if downloadData.Stat != "ok" {
		return nil, fmt.Errorf("DownloadPhoto is not ok\n%v\n", downloadData)
	}

	return &downloadData, nil
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type InfoPhotoData struct {
	Stat           string `pagser:"rsp->attr(stat)"`
	ID             string `pagser:"photo->attr(id)"`
	Secret         string `pagser:"photo->attr(secret)"`
	OriginalSecret string `pagser:"photo->attr(originalsecret)"`
	OriginalFormat string `pagser:"photo->attr(originalformat)"`
	Title          string `pagser:"title"`
	Description    string `pagser:"description"`
	Tags           []Tag  `pagser:"tag"`
}

type Tag struct {
	Name string `pagser:"->text()"`
}

func infoPhoto(parser *pagser.Pagser, photo PhotoFlickr) (*InfoPhotoData, error) {
	r := &Request{
		Host: "https://api.flickr.com/services/rest/?",
		Args: map[string]string{
			"api_key":  utils.DotEnvVariable("FLICKR_PRIVATE_KEY"),
			"method":   "flickr.photos.getInfo",
			"photo_id": photo.ID,
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, err
	}

	var infoData InfoPhotoData
	err = parser.Parse(&infoData, string(body))
	if err != nil {
		return nil, err
	}
	fmt.Println(utils.ToJSON(infoData))

	if infoData.Stat != "ok" {
		return nil, fmt.Errorf("InfoPhoto is not ok\n%v\n", infoData)
	}
	if photo.ID != infoData.ID {
		return nil, fmt.Errorf("IDs do not match! search id: %s, info id: %s\n", photo.ID, infoData.ID)
	}
	if photo.Secret != infoData.Secret {
		return nil, fmt.Errorf("Secrets do not match for id: %s! search secret: %s, info secret: %s\n", photo.ID, photo.Secret, infoData.Secret)
	}
	return &infoData, nil
}
