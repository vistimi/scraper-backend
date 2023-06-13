package adapter

import (
	hostModel "scraper-backend/src/driver/host/model"

	"github.com/foolin/pagser"
	"github.com/hbagdi/go-unsplash/unsplash"
)

// TODO: Marshal and Unmarshal types from driver to adapter
type DriverApiFlickr interface {
	GetFile(url string) ([]byte, error)
	SearchPhotosPerPage(parser *pagser.Pagser, licenseID string, tags string, page string) (*hostModel.SearchPhotPerPageData, error)
	DownloadPhoto(parser *pagser.Pagser, id string) (*hostModel.DownloadPhotoData, error)
	InfoPhoto(parser *pagser.Pagser, photo hostModel.PhotoFlickr) (*hostModel.InfoPhotoData, error)
}

type DriverApiUnsplash interface {
	GetFile(url string) ([]byte, error)
	GetPerPage() int
	SearchPhotosPerPage(tag string, page int) (*unsplash.PhotoSearchResult, error)
}

type DriverApiPexels interface {
	GetFile(url string) ([]byte, error)
	SearchPhotosPerPage(tag string, page int) (*hostModel.SearchPhotoResponsePexels, error)
}
