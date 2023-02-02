package usecase

import (
	"context"
	controllerModel "scraper-backend/src/adapter/controller/model"
	"scraper-backend/src/driver/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type ControllerPicture interface {
	ReadPictures(ctx context.Context, state string, projection *expression.ProjectionBuilder, filter *expression.ConditionBuilder) ([]controllerModel.Picture, error)
	ReadPicture(ctx context.Context, state string, primaryKey string, sortKey model.UUID) (*controllerModel.Picture, error)
	ReadPictureFile(ctx context.Context, origin, name, extension string) ([]byte, error)
	CreatePicture(ctx context.Context, id model.UUID, picture controllerModel.Picture, buffer []byte) error
	DeletePicture(ctx context.Context, primaryKey string, sortKey model.UUID) error
	DeletePictureAndFile(ctx context.Context, primaryKey string, sortKey model.UUID, name string) error
	DeletePicturesAndFiles(ctx context.Context, pictures []controllerModel.Picture) error
	CreatePictureTag(ctx context.Context, primaryKey string, sortKey model.UUID, tagID model.UUID, tag controllerModel.PictureTag) error
	UpdatePictureTag(ctx context.Context, primaryKey string, sortKey model.UUID, tagID model.UUID, tag controllerModel.PictureTag) error
	DeletePictureTag(ctx context.Context, primaryKey string, sortKey model.UUID, tagID model.UUID) error
	UpdatePictureCrop(ctx context.Context, primaryKey string, sortKey model.UUID, name string, imageSizeID model.UUID, box controllerModel.Box) error
	CreatePictureCrop(ctx context.Context, primaryKey string, sortKey model.UUID, id model.UUID, imageSizeID model.UUID, box controllerModel.Box) error
	CreatePictureCopy(ctx context.Context, primaryKey string, sortKey model.UUID, id model.UUID) error
	UpdatePictureTransfer(ctx context.Context, primaryKey string, sortKey model.UUID, from, to string) error
	CreatePictureBlocked(ctx context.Context, primaryKey string, sortKey model.UUID) error
	DeletePictureBlocked(ctx context.Context, primaryKey string, sortKey model.UUID) error
}

type ControllerTag interface {
	CreateTag(ctx context.Context, tag controllerModel.Tag) error
	CreateTagBlocked(ctx context.Context, tag controllerModel.Tag) error
	DeleteTag(ctx context.Context, primaryKey string, sortKey model.UUID) error
	ReadTags(ctx context.Context, primaryKey string) ([]controllerModel.Tag, error)
}

type ControllerUser interface {
	ReadUser(ctx context.Context, primaryKey string, sortKey model.UUID) (*controllerModel.User, error)
	CreateUser(ctx context.Context, user controllerModel.User) error
	DeleteUser(ctx context.Context, primaryKey string, sortKey model.UUID) error
	ReadUsers(ctx context.Context) ([]controllerModel.User, error)
}

type ControllerFlickr interface {
	SearchPhotos(ctx context.Context, quality string) error
}

type ControllerUnsplash interface {
	SearchPhotos(ctx context.Context, quality string) error
}

type ControllerPexels interface {
	SearchPhotos(ctx context.Context, quality string) error
}
