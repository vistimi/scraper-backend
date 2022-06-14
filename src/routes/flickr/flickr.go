package flickr

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/foolin/pagser"

	"path/filepath"

	"golang.org/x/exp/slices"

	"scrapper/src/mongodb"
	"scrapper/src/types"
	"scrapper/src/utils"

	"github.com/jinzhu/copier"

	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"regexp"
	"strings"
)

// Find all the photos with specific quality and folder directory.
func SearchPhoto(quality string, mongoClient *mongo.Client) ([]primitive.ObjectID, error) {

	var insertedIds []primitive.ObjectID

	parser := pagser.New()

	// If path is already a directory, MkdirAll does nothing and returns nil
	folderDir := utils.DotEnvVariable("FLICKR_PATH")
	err := os.MkdirAll(folderDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	collectionFlickr := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("FLICKR_COLLECTION"))

	// unwanted tags
	collectionUnwantedTags := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	res, err := mongodb.FindTags(collectionUnwantedTags)
	if err != nil {
		message := fmt.Sprintf("FindTags Unwated has failed: \n%v", err)
		return nil, errors.New(message)
	}
	var unwantedTags []string
	for _, tag := range res {
		unwantedTags = append(unwantedTags, strings.ToLower(tag.Name))
	}
	sort.Strings(unwantedTags)

	// wanted tags
	collectionWantedTags := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	res, err = mongodb.FindTags(collectionWantedTags)
	if err != nil {
		message := fmt.Sprintf("FindTags Wanted has failed: \n%v", err)
		return nil, errors.New(message)
	}
	var wantedTags []string
	for _, tag := range res {
		wantedTags = append(wantedTags, strings.ToLower(tag.Name))
	}
	sort.Strings(wantedTags)

	for _, wantedTag := range wantedTags {

		// all the commercial use licenses
		// https://www.flickr.com/services/api/flickr.photos.licenses.getInfo.html
		var licenseIdsNames = map[string]string{
			"4":  "Attribution License",
			"5":  "Attribution-ShareAlike License",
			"7":  "No known copyright restrictions",
			"9":  "Public Domain Dedication (CC0)",
			"10": "Public Domain Mark",
		}
		licenseIds := [5]string{"4", "5", "7", "9", "10"}
		for _, licenseId := range licenseIds {

			// start with the first page
			page := 1
			pageData, err := SearchPhotoPerPage(parser, licenseId, wantedTag, strconv.FormatUint(uint64(page), 10))
			if err != nil {
				message := fmt.Sprintf("SearchPhotoPerPage has failed: \n%v", err)
				return nil, errors.New(message)
			}

			for page := page; page <= int(pageData.Pages); page++ {
				pageData, err := SearchPhotoPerPage(parser, licenseId, wantedTag, strconv.FormatUint(uint64(page), 10))
				if err != nil {
					message := fmt.Sprintf("SearchPhotoPerPage has failed: \n%v", err)
					return nil, errors.New(message)
				}
				for _, photo := range pageData.Photos {

					// look for existing image
					_, err := mongodb.FindImageId(collectionFlickr, photo.Id)
					if err != nil {
						return nil, err
					}

					// extract the photo informations
					infoData, err := InfoPhoto(parser, photo)
					if err != nil {
						message := fmt.Sprintf("InfoPhoto has failed: \n%v", err)
						return nil, errors.New(message)
					}

					// only keep images with wanted tags
					idx := slices.IndexFunc(infoData.Tags, func(photoTag Tag) bool {
						imageTag := strings.ToLower(photoTag.Name)
						regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, imageTag)
						idx := slices.IndexFunc(unwantedTags, func(unwantedTag string) bool { 
							matched, err := regexp.Match(regexpMatch, []byte(unwantedTag))
							if err != nil {
								return false
							}
							return matched
						})
						if idx == -1 {
							return false
						} else {
							return true
						}
					})
					if idx == -1 {
						continue 
					}

					// extract the photo download link
					downloadData, err := DownloadPhoto(parser, photo.Id)
					if err != nil {
						message := fmt.Sprintf("DownloadPhoto has failed: \n%v", err)
						return nil, errors.New(message)
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
						message := fmt.Sprintf("Cannot find label %s and its derivatives %s in SearchPhoto! id %s has available the following:\n%v\n", label, regexpMatch, photo.Id, utils.ToJson(downloadData))
						return nil, errors.New(message)
					}

					// download photo into folder and rename it <id>.<format>
					fileName := fmt.Sprintf("%s.%s", photo.Id, infoData.OriginalFormat)
					path := fmt.Sprintf(filepath.Join(folderDir, fileName))
					err = DownloadFile(downloadData.Photos[idx].Source, path)
					if err != nil {
						return nil, err
					}

					var tags []types.Tag
					copier.Copy(&tags, &infoData.Tags)

					for i := 0; i < len(tags); i++ {
						tag := &tags[i]
						tag.Name = strings.ToLower(tag.Name)
						tag.CreationDate = time.Now()
						tag.Origin = "flickr"
					}

					document := types.Image{
						FlickrId:     photo.Id,
						Path:         path,
						Width:        downloadData.Photos[idx].Width,
						Height:       downloadData.Photos[idx].Height,
						Title:        infoData.Title,
						Description:  infoData.Description,
						License:      licenseIdsNames[licenseId],
						Tags:         tags,
						CreationDate: time.Now(),
					}

					insertedId, err := mongodb.InsertImage(collectionFlickr, document)
					if err != nil {
						return nil, err
					}
					insertedIds = append(insertedIds, insertedId)
				}
			}
		}
	}
	return insertedIds, nil
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type SearchPhotPerPageData struct {
	Stat    string  `pagser:"rsp->attr(stat)"`
	Page    uint    `pagser:"photos->attr(page)"`
	Pages   uint    `pagser:"photos->attr(pages)"`
	PerPage uint    `pagser:"photos->attr(perpage)"`
	Total   uint    `pagser:"photos->attr(total)"`
	Photos  []Photo `pagser:"photo"`
}
type Photo struct {
	Id     string `pagser:"->attr(id)"`
	Secret string `pagser:"->attr(secret)"`
	Title  string `pagser:"->attr(title)"`
}

// Search images for one page of max 500 images
func SearchPhotoPerPage(parser *pagser.Pagser, ids string, tags string, page string) (*SearchPhotPerPageData, error) {
	r := &Request{
		ApiKey: utils.DotEnvVariable("PRIVATE_KEY"),
		Method: "flickr.photos.search",
		Args: map[string]string{
			"tags":    tags,
			"license": ids,
			"media":   "photos",
			"page":    page,
		},
	}

	r.Sign(utils.DotEnvVariable("PUBLIC_KEY"))

	// log.Println(r.URL())

	response, err := r.Execute()
	if err != nil {
		return nil, err
	}

	var pageData SearchPhotPerPageData
	err = parser.Parse(&pageData, response)
	if err != nil {
		return nil, err
	}
	if pageData.Stat != "ok" {
		message := fmt.Sprintf("SearchPhotoPerPageRequest is not ok\n%v\n", utils.ToJson(pageData))
		return nil, errors.New(message)
	}
	if pageData.Page == 0 || pageData.Pages == 0 || pageData.PerPage == 0 || pageData.Total == 0 {
		message := fmt.Sprintf("Some informations are missing from SearchPhotoPerPage!")
		return nil, errors.New(message)
	}
	return &pageData, nil
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type DownloadPhotoSingleData struct {
	Label  string `pagser:"->attr(label)"`
	Width  uint   `pagser:"->attr(width)"`
	Height uint   `pagser:"->attr(height)"`
	Source string `pagser:"->attr(source)"`
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type DownloadPhotoData struct {
	Stat   string                    `pagser:"rsp->attr(stat)"`
	Photos []DownloadPhotoSingleData `pagser:"size"`
}

func DownloadPhoto(parser *pagser.Pagser, id string) (*DownloadPhotoData, error) {
	r := &Request{
		ApiKey: utils.DotEnvVariable("PRIVATE_KEY"),
		Method: "flickr.photos.getSizes",
		Args: map[string]string{
			"photo_id": id,
		},
	}

	r.Sign(utils.DotEnvVariable("PUBLIC_KEY"))

	// log.Println(r.URL())

	response, err := r.Execute()
	if err != nil {
		log.Fatalf("DownloadPhoto has failed: \n%v", err)
	}

	var downloadData DownloadPhotoData
	err = parser.Parse(&downloadData, response)
	if err != nil {
		log.Fatal(err)
	}
	if downloadData.Stat != "ok" {
		message := fmt.Sprintf("DownloadPhoto is not ok\n%v\n", utils.ToJson(downloadData))
		return nil, errors.New(message)
	}

	return &downloadData, nil
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type InfoPhotoData struct {
	Stat           string `pagser:"rsp->attr(stat)"`
	Id             string `pagser:"photo->attr(id)"`
	Secret         string `pagser:"photo->attr(secret)"`
	OriginalSecret string `pagser:"photo->attr(originalsecret)"`
	OriginalFormat string `pagser:"photo->attr(originalformat)"`
	Title          string `pagser:"title"`
	Description    string `pagser:"description"`
	Tags           []Tag `pagser:"tag"`
}

type Tag struct {
	Name string `pagser:"->text()"`
}

func InfoPhoto(parser *pagser.Pagser, photo Photo) (*InfoPhotoData, error) {
	r := &Request{
		ApiKey: utils.DotEnvVariable("PRIVATE_KEY"),
		Method: "flickr.photos.getInfo",
		Args: map[string]string{
			"photo_id": photo.Id,
		},
	}

	r.Sign(utils.DotEnvVariable("PUBLIC_KEY"))

	// log.Println(r.URL())

	response, err := r.Execute()
	if err != nil {
		return nil, err
	}

	var infoData InfoPhotoData
	err = parser.Parse(&infoData, response)
	if err != nil {
		return nil, err
	}
	if infoData.Stat != "ok" {
		message := fmt.Sprintf("InfoPhoto is not ok\n%v\n", utils.ToJson(infoData))
		return nil, errors.New(message)
	}
	if photo.Id != infoData.Id {
		message := fmt.Sprintf("IDs do not match! search id: %s, info id: %s\n", photo.Id, infoData.Id)
		return nil, errors.New(message)
	}
	if photo.Secret != infoData.Secret {
		message := fmt.Sprintf("Secrets do not match for id: %s! search secret: %s, info secret: %s\n", photo.Id, photo.Secret, infoData.Secret)
		return nil, errors.New(message)
	}
	return &infoData, nil
}
