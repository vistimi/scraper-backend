package gin

import (
	"context"
	"fmt"
	"time"
)

type ParamsSearchPhotoFlickr struct {
	Quality string `uri:"quality" binding:"required"`
}

func (d DriverServerGin) SearchPhotosFlickr(ctx context.Context, params ParamsSearchPhotoFlickr) (string, error) {
	if err := d.ControllerFlickr.SearchPhotos(ctx, params.Quality); err != nil {
		return "error", err
	}
	return "ok", nil
}

type ParamsSearchPhotoPexels struct {
	Quality string `uri:"quality" binding:"required"`
}

func (d DriverServerGin) SearchPhotosPexels(ctx context.Context, params ParamsSearchPhotoPexels) (string, error) {
	if err := d.ControllerPexels.SearchPhotos(ctx, params.Quality); err != nil {
		return "error", err
	}
	return "ok", nil
}

type ParamsSearchPhotoUnsplash struct {
	Quality    string `uri:"quality" binding:"required"`
	ImageStart int    `uri:"image_start"`
	ImageEnd   int    `uri:"image_end" binding:"required"`
}

func (d DriverServerGin) SearchPhotosUnsplash(ctx context.Context, params ParamsSearchPhotoUnsplash) ([]string, error) {
	t1 := time.Now()
	originIDs, err := d.ControllerUnsplash.SearchPhotos(ctx, params.Quality, params.ImageStart, params.ImageEnd)
	t2 := time.Now()
	fmt.Println("time: ", t2.Sub(t1))
	if err != nil {
		return nil, err
	}
	return originIDs, nil
}
