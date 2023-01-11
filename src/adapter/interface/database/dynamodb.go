package adapter

import (
	"context"
	controllerModel "scraper-backend/src/adapter/controller/model"
	dynamodbModel "scraper-backend/src/driver/database/dynamodb/model"

	"github.com/google/uuid"
)

type DriverDynamodb interface {
	ReadPicture(ctx context.Context, primaryKey, sortKey string) (*dynamodbModel.Picture, error)
	ReadPictures(ctx context.Context, primaryKey string) ([]dynamodbModel.Picture, error)
	ReadPicturesA(ctx context.Context, primaryKey, attributes string) ([]dynamodbModel.Picture, error)
	CreatePicture(ctx context.Context, picture controllerModel.Picture) error
	DeletePicture(ctx context.Context, primaryKey, sortKey string) error
	DeletePictureTag(ctx context.Context, primaryKey, sortKey string, tagID uuid.UUID) error
	CreatePictureTag(ctx context.Context, primaryKey, sortKey string, tag dynamodbModel.PictureTag) error
	UpdatePictureTag(ctx context.Context, primaryKey, sortKey string, tag map[uuid.UUID]dynamodbModel.PictureTag) error
	UpdatePictureSize(ctx context.Context, primaryKey, sortKey string, sizeMap map[uuid.UUID]dynamodbModel.PictureSize) error
}
