package routes

import (
	"errors"
	"fmt"
	"scraper/src/mongodb"
	"scraper/src/types"
	"scraper/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ParamsFindImagesIDs struct {
	Origin string `uri:"origin" binding:"required"`
}

func FindImagesIDs(mongoClient *mongo.Client, params ParamsFindImagesIDs) ([]types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	query := bson.M{"origin": params.Origin}
	options := options.Find().SetProjection(bson.M{"_id": 1})
	return mongodb.FindMany[types.Image](collectionImages, query, options)
}

type ParamsFindImage struct {
	ID string `uri:"id" binding:"required"`
}

func FindImage(mongoClient *mongo.Client, params ParamsFindImage) (*types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
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
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	return mongodb.FindOne[types.Image](collectionImagesUnwanted, bson.M{"origin": body.Origin, "originID": body.OriginID})
}

func FindImagesUnwanted(mongoClient *mongo.Client) ([]types.Image, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	return mongodb.FindMany[types.Image](collectionImagesUnwanted, bson.M{})
}

// Body for the RemoveImage request
type BodyRemoveImage struct {
	Origin string // image origin
	ID     primitive.ObjectID
}

func RemoveImage(mongoClient *mongo.Client, body BodyRemoveImage) (*int64, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	return mongodb.RemoveImageAndFile(collectionImages, body.ID, body.Origin)
}

func RemoveImageUnwanted(mongoClient *mongo.Client, body BodyRemoveImage) (*int64, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	return mongodb.RemoveImage(collectionImagesUnwanted, body.ID, body.Origin)
}

func UpdateImageTagsPush(mongoClient *mongo.Client, body types.BodyUpdateImageTagsPush) (*types.Image, error) {
	if body.ID == primitive.NilObjectID {
		return nil, errors.New("Body not valid, ID empty")
	}
	for _, tag := range body.Tags {
		if tag.Origin.Box.X == nil || tag.Origin.Box.Y == nil || tag.Origin.Box.Width == nil || tag.Origin.Box.Height == nil {
			return nil, fmt.Errorf("Body not valid, box fields missing: x=%d, y=%d, w=%d, h=%d", *tag.Origin.Box.X, *tag.Origin.Box.Y, *tag.Origin.Box.Width, *tag.Origin.Box.Height)
		}
	}
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	return mongodb.UpdateImageTagsPush(collectionImages, body)
}

func UpdateImageTagsPull(mongoClient *mongo.Client, body types.BodyUpdateImageTagsPull) (interface{}, error) {
	if body.ID == primitive.NilObjectID {
		return nil, errors.New("Body not valid, ID empty")
	}
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	return mongodb.UpdateImageTagsPull(collectionImages, body)
}

func UpdateImageFile(mongoClient *mongo.Client, body types.BodyUpdateImageFile) (*types.Image, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	return mongodb.UpdateImageFile(collectionImages, body)
}
