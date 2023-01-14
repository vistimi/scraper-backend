package controller

import (
	"context"
	controllerModel "scraper-backend/src/adapter/controller/model"
	databaseInterface "scraper-backend/src/adapter/interface/database"
)

type controllerUser struct {
	Dynamodb databaseInterface.DriverDynamodbUser
}

func (c controllerUser) CreateUser(ctx context.Context, user controllerModel.User) error {
	return c.Dynamodb.CreateUser(ctx, user)
}

func (c controllerUser) DeleteUser(ctx context.Context, user controllerModel.User) error {
	return c.Dynamodb.DeleteUser(ctx, user.Origin, user.Name)
}

func (c controllerUser) ReadUsers(ctx context.Context, user controllerModel.User) ([]controllerModel.User, error) {
	return c.Dynamodb.ReadUsers(ctx, user.Origin)
}
