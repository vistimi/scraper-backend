package usecase

import (
	"context"
	controllerModel "scraper-backend/src/adapter/controller/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
)

type ControllerPicture interface {
	ReadPictures(ctx context.Context, state string, projection *expression.ProjectionBuilder, filter *expression.ConditionBuilder) ([]controllerModel.Picture, error)
	ReadPicture(ctx context.Context, state string, primaryKey string, sortKey uuid.UUID) (*controllerModel.Picture, error)
	ReadPictureFile(ctx context.Context, origin, name, extension string) ([]byte, error)
	CreatePicture(ctx context.Context, id uuid.UUID, picture controllerModel.Picture) error
	DeletePicture(ctx context.Context, primaryKey string, sortKey uuid.UUID) error
	DeletePictureAndFile(ctx context.Context, primaryKey string, sortKey uuid.UUID, name string) error
	DeletePicturesAndFiles(ctx context.Context, pictures []controllerModel.Picture) error
	CreatePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID, tag controllerModel.PictureTag) error
	UpdatePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID, tag controllerModel.PictureTag) error
	DeletePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID) error
	UpdatePictureCrop(ctx context.Context, primaryKey string, sortKey uuid.UUID, name string, imageSizeID uuid.UUID, box controllerModel.Box) error
	CreatePictureCrop(ctx context.Context, primaryKey string, sortKey uuid.UUID, id uuid.UUID, imageSizeID uuid.UUID, box controllerModel.Box) error
	CreatePictureCopy(ctx context.Context, primaryKey string, sortKey uuid.UUID, id uuid.UUID) error
	UpdatePictureTransfer(ctx context.Context, primaryKey string, sortKey uuid.UUID, from, to string) error
	CreatePictureBlocked(ctx context.Context, primaryKey string, sortKey uuid.UUID) error
	DeletePictureBlocked(ctx context.Context, primaryKey string, sortKey uuid.UUID) error
}

type ControllerTag interface {
	CreateTag(ctx context.Context, tag controllerModel.Tag) error
	CreateTagBlocked(ctx context.Context, tag controllerModel.Tag) error
	DeleteTag(ctx context.Context, primaryKey string, sortKey uuid.UUID) error
	ReadTags(ctx context.Context, primaryKey string) ([]controllerModel.Tag, error)
}

type ControllerUser interface {
	CreateUser(ctx context.Context, user controllerModel.User) error
	DeleteUser(ctx context.Context, primaryKey string, sortKey uuid.UUID) error
	ReadUsers(ctx context.Context) ([]controllerModel.User, error)
}
