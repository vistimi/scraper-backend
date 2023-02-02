package controller

import (
	"context"
	controllerModel "scraper-backend/src/adapter/controller/model"
	databaseInterface "scraper-backend/src/driver/interface/database"
	model "scraper-backend/src/driver/model"
)

type ControllerUser struct {
	Dynamodb databaseInterface.DriverDynamodbUser
}

func (c ControllerUser) CreateUser(ctx context.Context, user controllerModel.User) error {
	return c.Dynamodb.CreateUser(ctx, user)
}

func (c ControllerUser) DeleteUser(ctx context.Context, primaryKey string, sortKey model.UUID) error {
	return c.Dynamodb.DeleteUser(ctx, primaryKey, sortKey)
}

func (c ControllerUser) ReadUsers(ctx context.Context) ([]controllerModel.User, error) {
	return c.Dynamodb.ScanUsers(ctx)
}

func (c ControllerUser) ReadUser(ctx context.Context, primaryKey string, sortKey model.UUID) (*controllerModel.User, error) {
	return c.Dynamodb.ReadUser(ctx, primaryKey, sortKey)
}
