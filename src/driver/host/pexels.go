package host

import (
	"encoding/json"
	"fmt"
	"scraper-backend/src/util"

	hostModel "scraper-backend/src/driver/host/model"
)

type DriverApiPexels struct {
}

func (d *DriverApiPexels) SearchPhotosPerPage(tag string, page int) (*hostModel.SearchPhotoResponsePexels, error) {
	r := &Request{
		Host: "https://api.pexels.com/v1/search?",
		Args: map[string]string{
			"query":    tag,
			"per_page": "80", // default 15, max 80
			"page":     fmt.Sprint(page),
		},
		Header: map[string][]string{
			"Authorization": {util.GetEnvVariable("PEXELS_PUBLIC_KEY")},
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, err
	}

	var searchPerPage hostModel.SearchPhotoResponsePexels
	err = json.Unmarshal(body, &searchPerPage)
	if err != nil {
		return nil, err
	}
	return &searchPerPage, nil
}

func (d *DriverApiPexels) GetFile(url string) ([]byte, error) {
	return GetFile(url)
}