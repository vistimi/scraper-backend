package gin

import "context"

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
	Quality string `uri:"quality" binding:"required"`
}

func (d DriverServerGin) SearchPhotosUnsplash(ctx context.Context, params ParamsSearchPhotoUnsplash) (string, error) {
	if err := d.ControllerUnsplash.SearchPhotos(ctx, params.Quality); err != nil {
		return "error", err
	}
	return "ok", nil
}
