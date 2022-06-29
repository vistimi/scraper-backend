package routes

import (
	"errors"
	"scrapper/src/mongodb"
	"scrapper/src/types"
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ParamsFindImagesIDs struct {
	Origin string `uri:"origin" binding:"required"`
}

func FindImagesIDs(mongoClient *mongo.Client, params ParamsFindImagesIDs) ([]types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	query := bson.M{"origin": params.Origin}
	options := options.Find().SetProjection(bson.M{"_id": 1})
	return mongodb.FindMany[types.Image](collectionImages, query, options)
}

type ParamsFindImage struct {
	ID string `uri:"id" binding:"required"`
}

func FindImage(mongoClient *mongo.Client, params ParamsFindImage) (*types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	imageID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.FindOne[types.Image](collectionImages, bson.M{"_id": imageID})
}

type BodyFindImageUnwanted struct {
	Origin   string `bson:"origin" json:"origin"`
	OriginID string `bson:"originID" json:"originID"`
}

func FindImageUnwanted(mongoClient *mongo.Client, body BodyFindImageUnwanted) (*types.Image, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	return mongodb.FindOne[types.Image](collectionImagesUnwanted, bson.M{"origin": body.Origin, "originID": body.OriginID})
}

func FindImagesUnwanted(mongoClient *mongo.Client) ([]types.Image, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	return mongodb.FindMany[types.Image](collectionImagesUnwanted, bson.M{})
}

// Body for the RemoveImage request
type BodyRemoveImage struct {
	Origin string // image origin
	ID     primitive.ObjectID
}

func RemoveImage(mongoClient *mongo.Client, body BodyRemoveImage) (*int64, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	return mongodb.RemoveImageAndFile(collectionImages, body.ID, body.Origin)
}

func RemoveImageUnwanted(mongoClient *mongo.Client, body BodyRemoveImage) (*int64, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	return mongodb.RemoveImage(collectionImagesUnwanted, body.ID, body.Origin)
}

func UpdateImage(mongoClient *mongo.Client, body types.BodyUpdateImage) (*types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	if body.ID == primitive.NilObjectID {
		return nil, errors.New("Body not valid, ID empty")
	}
	return mongodb.UpdateImage(collectionImages, body)
}
