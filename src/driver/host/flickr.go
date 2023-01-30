package host

import (
	"errors"
	"fmt"
	"scraper-backend/src/util"

	"github.com/foolin/pagser"
	hostModel "scraper-backend/src/driver/host/model"
)

type DriverApiFlickr struct {
}

// Search images for one page of max 500 images
func (d *DriverApiFlickr) SearchPhotosPerPage(parser *pagser.Pagser, licenseID string, tags string, page string) (*hostModel.SearchPhotPerPageData, error) {
	r := &Request{
		Host: "https://api.flickr.com/services/rest/?",
		Args: map[string]string{
			"api_key":  util.GetEnvVariable("FLICKR_PUBLIC_KEY"),
			"method":   "flickr.photos.search",
			"tags":     tags,
			"license":  licenseID,
			"media":    "photos",
			"per_page": "500", // 100 default, max 500
			"page":     page,
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, err
	}

	var pageData hostModel.SearchPhotPerPageData
	err = parser.Parse(&pageData, string(body))
	if err != nil {
		return nil, err
	}
	if pageData.Stat != "ok" {
		return nil, fmt.Errorf("SearchPhotoPerPageRequest is not ok%v", pageData)
	}
	if pageData.Page == 0 || pageData.Pages == 0 || pageData.PerPage == 0 || pageData.Total == 0 {
		return nil, errors.New("some informations are missing from SearchPhotoPerPage")
	}
	return &pageData, nil
}

func (d *DriverApiFlickr) DownloadPhoto(parser *pagser.Pagser, id string) (*hostModel.DownloadPhotoData, error) {
	r := &Request{
		Host: "https://api.flickr.com/services/rest/?",
		Args: map[string]string{
			"api_key":  util.GetEnvVariable("FLICKR_PUBLIC_KEY"),
			"method":   "flickr.photos.getSizes",
			"photo_id": id,
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, fmt.Errorf("DownloadPhoto has failed: %v", err)
	}

	var downloadData hostModel.DownloadPhotoData
	err = parser.Parse(&downloadData, string(body))
	if err != nil {
		return nil, err
	}

	if downloadData.Stat != "ok" {
		return nil, fmt.Errorf("DownloadPhoto is not ok%v", downloadData)
	}

	return &downloadData, nil
}

func (d *DriverApiFlickr) InfoPhoto(parser *pagser.Pagser, photo hostModel.PhotoFlickr) (*hostModel.InfoPhotoData, error) {
	r := &Request{
		Host: "https://api.flickr.com/services/rest/?",
		Args: map[string]string{
			"api_key":  util.GetEnvVariable("FLICKR_PUBLIC_KEY"),
			"method":   "flickr.photos.getInfo",
			"photo_id": photo.ID,
		},
	}
	// fmt.Println(r.URL())

	body, err := r.ExecuteGET()
	if err != nil {
		return nil, err
	}

	var infoData hostModel.InfoPhotoData
	err = parser.Parse(&infoData, string(body))
	if err != nil {
		return nil, err
	}

	if infoData.Stat != "ok" {
		return nil, fmt.Errorf("InfoPhoto is not ok%v", infoData)
	}
	if photo.ID != infoData.ID {
		return nil, fmt.Errorf("IDs do not match! search id: %s, info id: %s", photo.ID, infoData.ID)
	}
	if photo.Secret != infoData.Secret {
		return nil, fmt.Errorf("secrets do not match for id: %s! search secret: %s, info secret: %s", photo.ID, photo.Secret, infoData.Secret)
	}
	return &infoData, nil
}

func (d *DriverApiFlickr) GetFile(url string) ([]byte, error) {
	return GetFile(url)
}
