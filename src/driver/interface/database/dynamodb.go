package adapter

import (
	"context"
	controllerModel "scraper-backend/src/adapter/controller/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
)

type DriverDynamodbPicture interface {
	ReadPicture(ctx context.Context, primaryKey string, sortKey uuid.UUID) (*controllerModel.Picture, error)
	ReadPictures(ctx context.Context, projection *expression.ProjectionBuilder, filter *expression.ConditionBuilder) ([]controllerModel.Picture, error)
	CreatePicture(ctx context.Context, id uuid.UUID, picture controllerModel.Picture) error
	DeletePicture(ctx context.Context, primaryKey string, sortKey uuid.UUID) error
	DeletePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID) error
	CreatePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID, tag controllerModel.PictureTag) error
	UpdatePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID, tag controllerModel.PictureTag) error
	CreatePictureSize(ctx context.Context, primaryKey string, sortKey uuid.UUID, size controllerModel.PictureSize) error
}

type DriverDynamodbTag interface {
	ReadTags(ctx context.Context, primaryKey string) ([]controllerModel.Tag, error)
	CreateTag(ctx context.Context, picture controllerModel.Tag) error
	DeleteTag(ctx context.Context, primaryKey string, sortKey uuid.UUID) error
	ScanTags(ctx context.Context) ([]controllerModel.Tag, error)
}

type DriverDynamodbUser interface {
	ReadUser(ctx context.Context, primaryKey string, sortKey uuid.UUID) (*controllerModel.User, error)
	ReadUsers(ctx context.Context, primaryKey string) ([]controllerModel.User, error)
	ScanUsers(ctx context.Context) ([]controllerModel.User, error)
	CreateUser(ctx context.Context, picture controllerModel.User) error
	DeleteUser(ctx context.Context, primaryKey string, sortKey uuid.UUID) error
}
