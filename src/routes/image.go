package routes

import (
	"errors"
	"scrapper/src/mongodb"
	"scrapper/src/types"
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ParamsFindImagesIds struct {
	Collection string `uri:"collection" binding:"required"`
}

func FindImagesIds(mongoClient *mongo.Client, params ParamsFindImagesIds) ([]types.Image, error) {
	var collection *mongo.Collection
	switch params.Collection {
	case "flickr":
		collection = mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("FLICKR_COLLECTION"))
	default:
		return nil, errors.New("Params not valid, you must give a correct collection!")
	}
	return mongodb.FindImagesIds(collection)
}


func UpdateImage(mongoClient *mongo.Client, body types.BodyUpdateImage) (*types.Image, error) {
	var collection *mongo.Collection
	switch body.Collection {
	case "flickr":
		collection = mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("FLICKR_COLLECTION"))
	default:
		return nil, errors.New("Body not valid, you must give a correct collection!")
	}
	if body.Id == primitive.NilObjectID {
		return nil, errors.New("Body not valid, Id empty")
	}
	return mongodb.UpdateImage(collection, body)
}
