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
	Origin string `uri:"origin" binding:"required"`
}

func FindImagesIDs(mongoClient *mongo.Client, params ParamsFindImagesIDs) ([]types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	return mongodb.FindImagesIDs(collectionImages, bson.M{"origin": params.Origin})
}

type ParamsFindImage struct {
	ID         string `uri:"id" binding:"required"`
}

func FindImage(mongoClient *mongo.Client, params ParamsFindImage) (*types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	imageID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.FindImageByID(collectionImages, imageID)
}

// Body for the RemoveImage request
type BodyRemoveImage struct {
	Origin string		// image origin
	ID         primitive.ObjectID
}

func RemoveImage(mongoClient *mongo.Client, body BodyRemoveImage) (*int64, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	return mongodb.RemoveImageAndFile(collectionImages, body.Origin, body.ID)
}

func UpdateImage(mongoClient *mongo.Client, body types.BodyUpdateImage) (*types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	if body.ID == primitive.NilObjectID {
		return nil, errors.New("Body not valid, ID empty")
	}
	return mongodb.UpdateImage(collectionImages, body)
}
