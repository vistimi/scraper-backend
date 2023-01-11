package controller

import (
	"context"
	"path/filepath"
	controllerModel "scraper-backend/src/adapter/controller/model"
	databaseInterface "scraper-backend/src/adapter/interface/database"
	storageInterface "scraper-backend/src/adapter/interface/storage"
)

type ControllerPicture struct {
	S3         storageInterface.DriverS3
	BucketName string
	Dynamodb   databaseInterface.DriverDynamodb
	PrimaryKey string
	SortKey    string
}

func (c *ControllerPicture) CreatePicture(ctx context.Context, picture controllerModel.Picture) error {
	return c.Dynamodb.CreatePicture(ctx, picture)
}

func (c *ControllerPicture) DeletePicture(ctx context.Context, picture controllerModel.Picture) error {
	return c.Dynamodb.DeletePicture(ctx, c.PrimaryKey, c.SortKey)
}

func (c *ControllerPicture) DeletePictureAndFile(ctx context.Context, picture controllerModel.Picture) error {
	picture := c.Dynamodb.ReadPictureA()
	if err := c.Dynamodb.DeletePicture(ctx, c.PrimaryKey, c.SortKey); err != nil {
		return err
	}
	path := filepath.Join(picture.Origin, picture.Name)
	return c.S3.ItemDelete(ctx, c.BucketName, path)
}
