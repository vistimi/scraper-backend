// https://www.flickr.com/services/api/
// machinetags

package flickr

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/foolin/pagser"

	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and return the value of the key
func DotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func toJson(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}

type SearchPhotoData struct {
	Stat    string `pagser:"rsp->attr(stat)"`
	Page    uint   `pagser:"photos->attr(page)"`
	Pages   uint   `pagser:"photos->attr(pages)"`
	PerPage uint   `pagser:"photos->attr(perpage)"`
	Total   uint   `pagser:"photos->attr(total)"`
	Photos  []struct {
		ID    uint    `pagser:"->attr(id)"`
		Owner string `pagser:"->attr(owner)"`
		Secret string `pagser:"->attr(secret)"`
		Title string `pagser:"->attr(title)"`
	} `pagser:"photo"`
}

func SearchPhoto(ids string, tags string) {

	page := 1

	response := SearchPhotoPerPage(ids, tags, strconv.FormatUint(uint64(page), 10))
	p := pagser.New()

	var pageData SearchPhotoData
	err := p.Parse(&pageData, response)
	if err != nil {
		log.Fatal(err)
	}

	if pageData.Stat != "ok" || pageData.Page == 0 || pageData.Pages == 0 || pageData.PerPage == 0 || pageData.Total == 0 {return}

	for _, photo := range pageData.Photos {
		response := DownloadPhoto(photo.ID)
		var downloadData DownloadPhotoData
		err := p.Parse(&downloadData, response)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Page data json: \n-------------\n%v\n-------------\n", toJson(downloadData))
		break
	}

	//print data
	// log.Printf("Page data json: \n-------------\n%v\n-------------\n", toJson(data))
}

// id: "4, 5, 7, 9, 10".
// Find all the photos with specific tags and licenses.
// Example: https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
func SearchPhotoPerPage(ids string, tags string, page string) string {
	r := &Request{
		ApiKey: DotEnvVariable("PRIVATE_KEY"),
		Method: "flickr.photos.search",
		Args: map[string]string{
			"tags":    tags,
			"license": ids,
			"media":   "photos",
			"page": page,
		},
	}

	r.Sign(DotEnvVariable("PUBLIC_KEY"))

	log.Println(r.URL())

	response, err := r.Execute()
	if err != nil {
		log.Println(err)
	}

	return response
}

type DownloadPhotoData struct {
	Stat    string `pagser:"rsp->attr(stat)"`
	Photos  []struct {
		Label    string    `pagser:"->attr(label)"`
		Width uint `pagser:"->attr(width)"`
		Height uint `pagser:"->attr(height)"`
		Source string `pagser:"->attr(source)"`
	} `pagser:"size"`
}

func DownloadPhoto(id uint) string {
	r := &Request{
		ApiKey: DotEnvVariable("PRIVATE_KEY"),
		Method: "flickr.photos.getSizes",
		Args: map[string]string{
			"photo_id":    strconv.FormatUint(uint64(id), 10),
		},
	}

	r.Sign(DotEnvVariable("PUBLIC_KEY"))

	log.Println(r.URL())

	response, err := r.Execute()
	if err != nil {
		log.Println(err)
	}

	return response
}
