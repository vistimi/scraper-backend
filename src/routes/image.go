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

type ParamsFindImagesIds struct {
	Collection string `uri:"collection" binding:"required"`
}

func FindImagesIds(mongoClient *mongo.Client, params ParamsFindImagesIds) ([]types.Image, error) {
	collection, err := utils.ImageCollectionSelection(mongoClient, params.Collection)
	if err != nil {
		return nil, err
	}
	return mongodb.FindImagesIds(collection, bson.M{})
}

type ParamsFindImage struct {
	Collection string `uri:"collection" binding:"required"`
	Id         string `uri:"id" binding:"required"`
}

func FindImage(mongoClient *mongo.Client, params ParamsFindImage) (*types.Image, error) {
	collection, err := utils.ImageCollectionSelection(mongoClient, params.Collection)
	if err != nil {
		return nil, err
	}
	imageId, err := primitive.ObjectIDFromHex(params.Id)
	if err != nil {
		return nil, err
	}
	return mongodb.FindImage(collection, imageId)
}

// Body for the RemoveImage request
type BodyRemoveImage struct {
	Collection string		// image collection
	Id         primitive.ObjectID
}

func RemoveImage(mongoClient *mongo.Client, body BodyRemoveImage) (*int64, error) {
	collection, err := utils.ImageCollectionSelection(mongoClient, body.Collection)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveImageAndFile(collection, body.Collection, body.Id)
}

func UpdateImage(mongoClient *mongo.Client, body types.BodyUpdateImage) (*types.Image, error) {
	collection, err := utils.ImageCollectionSelection(mongoClient, body.Collection)
	if err != nil {
		return nil, err
	}
	if body.ID == primitive.NilObjectID {
		return nil, errors.New("Body not valid, Id empty")
	}
	return mongodb.UpdateImage(collection, body)
}
