package routes

import (
	"errors"
	"scrapper/src/mongodb"
	"scrapper/src/types"
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type ParamsFindImagesIDs struct {
	Collection string `uri:"collection" binding:"required"`
}

func FindImagesIDs(mongoClient *mongo.Client, params ParamsFindImagesIDs) ([]types.Image, error) {
	collection, err := utils.ImageCollectionSelection(mongoClient, params.Collection)
	if err != nil {
		return nil, err
	}
	return mongodb.FindImagesIDs(collection, bson.M{})
}

type ParamsFindImage struct {
	Collection string `uri:"collection" binding:"required"`
	ID         string `uri:"id" binding:"required"`
}

func FindImage(mongoClient *mongo.Client, params ParamsFindImage) (*types.Image, error) {
	collection, err := utils.ImageCollectionSelection(mongoClient, params.Collection)
	if err != nil {
		return nil, err
	}
	imageID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.FindImageByID(collection, imageID)
}

// Body for the RemoveImage request
type BodyRemoveImage struct {
	Collection string		// image collection
	ID         primitive.ObjectID
}

func RemoveImage(mongoClient *mongo.Client, body BodyRemoveImage) (*int64, error) {
	collection, err := utils.ImageCollectionSelection(mongoClient, body.Collection)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveImageAndFile(collection, body.Collection, body.ID)
}

func UpdateImage(mongoClient *mongo.Client, body types.BodyUpdateImage) (*types.Image, error) {
	collection, err := utils.ImageCollectionSelection(mongoClient, body.Collection)
	if err != nil {
		return nil, err
	}
	if body.ID == primitive.NilObjectID {
		return nil, errors.New("Body not valid, ID empty")
	}
	return mongodb.UpdateImage(collection, body)
}
