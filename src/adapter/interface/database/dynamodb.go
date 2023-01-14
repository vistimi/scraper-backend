package adapter

import (
	"context"
	controllerModel "scraper-backend/src/adapter/controller/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
)

type DriverDynamodbPicture interface {
	ReadPicture(ctx context.Context, primaryKey, sortKey string) (*controllerModel.Picture, error)
	ReadPictures(ctx context.Context, primaryKey string, projection *expression.ProjectionBuilder, filter *expression.ConditionBuilder) ([]controllerModel.Picture, error)
	CreatePicture(ctx context.Context, picture controllerModel.Picture) error
	DeletePicture(ctx context.Context, primaryKey, sortKey string) error
	DeletePictureTag(ctx context.Context, primaryKey, sortKey string, tagID uuid.UUID) error
	CreatePictureTag(ctx context.Context, primaryKey, sortKey string, tag controllerModel.PictureTag) error
	UpdatePictureTag(ctx context.Context, primaryKey, sortKey string, tagID uuid.UUID, tag controllerModel.PictureTag) error
	CreatePictureSize(ctx context.Context, primaryKey, sortKey string, size controllerModel.PictureSize) error
}

type DriverDynamodbTag interface {
	ReadTag(ctx context.Context, primaryKey, sortKey string) (*controllerModel.Tag, error)
	ReadTags(ctx context.Context, primaryKey string) ([]controllerModel.Tag, error)
	CreateTag(ctx context.Context, picture controllerModel.Tag) error
	DeleteTag(ctx context.Context, primaryKey, sortKey string) error
	ScanTags(ctx context.Context) ([]controllerModel.Tag, error)
}

type DriverDynamodbUser interface {
	ReadUsers(ctx context.Context, primaryKey string) ([]controllerModel.User, error)
	CreateUser(ctx context.Context, picture controllerModel.User) error
	DeleteUser(ctx context.Context, primaryKey, sortKey string) error
}
