package host

import (
	"fmt"
	"scraper-backend/src/util"

	"github.com/hbagdi/go-unsplash/unsplash"

	"encoding/json"
)

type DriverApiUnsplash struct {
}

func (d *DriverApiUnsplash) SearchPhotosPerPage(tag string, page int) (*unsplash.PhotoSearchResult, error) {
	r := &Request{
		Host: "https://api.unsplash.com/search/photos/?",
		Args: map[string]string{
			"client_id": util.GetEnvVariable("UNSPLASH_PUBLIC_KEY"),
			"per_page":  "80", // default 10
			"page":      fmt.Sprint(page),
			"query":     tag,
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, err
	}

	var searchPerPage unsplash.PhotoSearchResult
	err = json.Unmarshal(body, &searchPerPage)
	if err != nil {
		return nil, err
	}
	return &searchPerPage, nil
}

func (d *DriverApiUnsplash) GetFile(url string) ([]byte, error) {
	return GetFile(url)
}