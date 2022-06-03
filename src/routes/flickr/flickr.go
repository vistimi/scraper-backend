// https://www.flickr.com/services/api/
// machinetags

package flickr

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/foolin/pagser"

	"path/filepath"

	"golang.org/x/exp/slices"

	"dressme-scrapper/src/mongodb"
	"dressme-scrapper/src/types"
	"dressme-scrapper/src/utils"

	"github.com/jinzhu/copier"
	"gopkg.in/mgo.v2"
)

// To transform the structured extracted data from html into json.
// log.Println(toJson(data))
func toJson(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}

// id: "4, 5, 7, 9, 10".
// Find all the photos with specific tags and licenses.
// Example: https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
func SearchPhoto(licenseId string, tags string, quality string, folderDir string, mongoSession *mgo.Session) {

	// If path is already a directory, MkdirAll does nothing and returns nil
	err := os.MkdirAll(folderDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	p := pagser.New()

	// TODO: fetch in db the existing licenseId and skip those photos

	page := 1
	// TODO: loop through all pages

	response := SearchPhotoPerPage(licenseId, tags, strconv.FormatUint(uint64(page), 10))

	var pageData SearchPhotPerPageData
	err = p.Parse(&pageData, response)
	if err != nil {
		log.Fatal(err)
	}
	if pageData.Stat != "ok" {
		log.Printf("SearchPhotoPerPageRequest is not ok\n%v\n", toJson(pageData))
		// TODO: continue
	}
	if pageData.Page == 0 || pageData.Pages == 0 || pageData.PerPage == 0 || pageData.Total == 0 {
		log.Println("Some informations are missing from SearchPhotoPerPage!")
		// TODO: continue
	}

	for i, photo := range pageData.Photos {

		// extract the photo download link
		response := DownloadPhoto(photo.ID)
		var downloadData DownloadPhotoData
		err := p.Parse(&downloadData, response)
		if err != nil {
			log.Fatal(err)
		}
		if pageData.Stat != "ok" {
			log.Printf("DownloadPhoto is not ok\n%v\n", toJson(downloadData))
			continue
		}

		// get the download link for the correct resolution
		label := "Medium"
		idx := slices.IndexFunc(downloadData.Photos, func(c DownloadPhotoSingleData) bool { return c.Label == label })
		if idx == -1 {
			// TODO: download higher resolution or Medium_...
			log.Printf("Cannot find label in SearchPhoto! Page %d, image %d, label %s, id %s\n", page, i, label, photo.ID)
			continue
		}

		// extract the photo informations
		response = InfoPhoto(photo.ID)
		var infoData InfoPhotoData
		err = p.Parse(&infoData, response)
		if err != nil {
			log.Fatal(err)
		}
		if pageData.Stat != "ok" {
			log.Printf("InfoPhoto is not ok\n%v\n", toJson(infoData))
			continue
		}

		if photo.ID != infoData.ID {
			log.Printf("IDs do not match! search id: %s, info id: %s\n", photo.ID, infoData.ID)
			continue
		}
		if photo.Secret != infoData.Secret {
			log.Printf("Secrets do not match for id: %s! search secret: %s, info secret: %s\n", photo.ID, photo.Secret, infoData.Secret)
			continue
		}

		// download photo into folder and rename it <id>.<format>
		fileName := fmt.Sprintf("%s.%s", photo.ID, infoData.OriginalFormat)
		path := fmt.Sprintf(filepath.Join(folderDir, fileName))
		err = downloadFile(downloadData.Photos[idx].Source, path)
		if err != nil {
			log.Fatal(err)
			// log.Printf("Cannot download photo id: %s into path: %s\n%v\n", photo.ID, path, err)
			continue
		}

		sessionCopy := mongoSession.Copy()
		defer sessionCopy.Close()
		collection := sessionCopy.DB(utils.DotEnvVariable("SCRAPPER_DB")).C(utils.DotEnvVariable("FLICKR_COLL"))

		var tags []types.FlickTag
		copier.Copy(&tags, &infoData.Tags)
		document := types.FlickrImage{
			FlickrId:    photo.ID,
			Path:        path,
			Width:       downloadData.Photos[idx].Width,
			Height:      downloadData.Photos[idx].Height,
			Title:       infoData.Title,
			Description: infoData.Description,
			License:     licenseId,
			Tags:        tags,
		}

		log.Println(document)
		return

		// return fmt.Sprint(inserted.InsertedID), nil
		err = mongodb.InsertImage(collection, document)
		if err != nil {
			log.Fatal(err)
			// log.Printf("Cannot download photo id: %s into path: %s\n%v\n", photo.ID, path, err)
			continue
		}
	}
}

type SearchPhotPerPageData struct {
	Stat    string `pagser:"rsp->attr(stat)"`
	Page    uint   `pagser:"photos->attr(page)"`
	Pages   uint   `pagser:"photos->attr(pages)"`
	PerPage uint   `pagser:"photos->attr(perpage)"`
	Total   uint   `pagser:"photos->attr(total)"`
	Photos  []struct {
		ID     string `pagser:"->attr(id)"`
		Secret string `pagser:"->attr(secret)"`
		Title  string `pagser:"->attr(title)"`
	} `pagser:"photo"`
}

// Search images for one page of max 500 images
func SearchPhotoPerPage(ids string, tags string, page string) string {
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

	log.Println(r.URL())

	response, err := r.Execute()
	if err != nil {
		log.Println(err)
	}

	return response
}

type DownloadPhotoSingleData struct {
	Label  string `pagser:"->attr(label)"`
	Width  uint   `pagser:"->attr(width)"`
	Height uint   `pagser:"->attr(height)"`
	Source string `pagser:"->attr(source)"`
}
type DownloadPhotoData struct {
	Stat   string                    `pagser:"rsp->attr(stat)"`
	Photos []DownloadPhotoSingleData `pagser:"size"`
}

func DownloadPhoto(id string) string {
	r := &Request{
		ApiKey: utils.DotEnvVariable("PRIVATE_KEY"),
		Method: "flickr.photos.getSizes",
		Args: map[string]string{
			"photo_id": id,
		},
	}

	r.Sign(utils.DotEnvVariable("PUBLIC_KEY"))

	log.Println(r.URL())

	response, err := r.Execute()
	if err != nil {
		log.Println(err)
	}

	return response
}

type InfoPhotoData struct {
	Stat           string `pagser:"rsp->attr(stat)"`
	ID             string `pagser:"photo->attr(id)"`
	Secret         string `pagser:"photo->attr(secret)"`
	OriginalSecret string `pagser:"photo->attr(originalsecret)"`
	OriginalFormat string `pagser:"photo->attr(originalformat)"`
	Title          string `pagser:"title"`
	Description    string `pagser:"description"`
	Tags           []struct {
		ID   string `pagser:"->attr(id)"`
		Name string `pagser:"->text()"`
	} `pagser:"tag"`
}

func InfoPhoto(id string) string {
	r := &Request{
		ApiKey: utils.DotEnvVariable("PRIVATE_KEY"),
		Method: "flickr.photos.getInfo",
		Args: map[string]string{
			"photo_id": id,
		},
	}

	r.Sign(utils.DotEnvVariable("PUBLIC_KEY"))

	log.Println(r.URL())

	response, err := r.Execute()
	if err != nil {
		log.Println(err)
	}

	return response
}
