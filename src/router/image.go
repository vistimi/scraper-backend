package router

import (
	"errors"
	"fmt"
	"path/filepath"
	"scraper/src/mongodb"
	"scraper/src/types"
	"scraper/src/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ParamsFindImageFile struct {
	Origin    string `uri:"origin" binding:"required"`
	Name      string `uri:"name" binding:"required"`
	Extension string `uri:"extension" binding:"required"`
}

func FindImageFile(s3Client *s3.Client, params ParamsFindImageFile) (*DataSchema, error) {
	fileName := fmt.Sprintf("%s.%s", params.Name, params.Extension)
	path := filepath.Join(params.Origin, fileName)

	buffer, err := utils.GetItemS3(s3Client, path)
	if err != nil {
		return nil, err
	}
	data := DataSchema{dataType: params.Extension, dataFile: buffer}
	return &data, nil
}

type ParamsFindImagesIDs struct {
	Origin     string `uri:"origin" binding:"required"`
	Collection string `uri:"collection" binding:"required"`
}

// FindImagesIDs get all the IDs of an image collection
func FindImagesIDs(mongoClient *mongo.Client, params ParamsFindImagesIDs) ([]types.Image, error) {
	collectionImages, err := utils.ImagesCollection(mongoClient, params.Collection)
	if err != nil {
		return nil, err
	}
	query := bson.M{"origin": params.Origin}
	options := options.Find().SetProjection(bson.M{"_id": 1})
	return mongodb.FindMany[types.Image](collectionImages, query, options)
}

type ParamsFindImage struct {
	ID         string `uri:"id" binding:"required"`
	Collection string `uri:"collection" binding:"required"`
}

// FindImage get a specific image
func FindImage(mongoClient *mongo.Client, params ParamsFindImage) (*types.Image, error) {
	collectionImages, err := utils.ImagesCollection(mongoClient, params.Collection)
	if err != nil {
		return nil, err
	}
	imageID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.FindOne[types.Image](collectionImages, bson.M{"_id": imageID})
}

// FindImagesUnwanted get all the unwanted images
func FindImagesUnwanted(mongoClient *mongo.Client) ([]types.Image, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	// no options needed because not much is stored for unwanted images
	return mongodb.FindMany[types.Image](collectionImagesUnwanted, bson.M{})
}

// Body for the RemoveImage request
type ParamsRemoveImage struct {
	ID string `uri:"id" binding:"required"`
}

// RemoveImageAndFile removes in db and file of a pending image
func RemoveImageAndFile(s3Client *s3.Client, mongoClient *mongo.Client, params ParamsRemoveImage) (*int64, error) {
	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("IMAGES_PENDING_COLLECTION"))
	imageID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveImageAndFile(s3Client, collectionImagesPending, imageID)
}

// RemoveImage removes in db an unwanted image
func RemoveImage(mongoClient *mongo.Client, params ParamsRemoveImage) (*int64, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	imageID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveImage(collectionImagesUnwanted, imageID)
}

// UpdateImageTagsPush add tags to a pending image
func UpdateImageTagsPush(mongoClient *mongo.Client, body types.BodyUpdateImageTagsPush) (*int64, error) {
	if body.ID == primitive.NilObjectID {
		return nil, errors.New("body not valid, ID empty")
	}
	for _, tag := range body.Tags {
		if tag.Origin.Box.Tlx == nil || tag.Origin.Box.Tly == nil || tag.Origin.Box.Width == nil || tag.Origin.Box.Height == nil {
			return nil, fmt.Errorf("body not valid, box fields missing: %v", tag.Origin.Box)
		}
	}
	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("IMAGES_PENDING_COLLECTION"))
	return mongodb.UpdateImageTagsPush(collectionImagesPending, body)
}

// UpdateImageTagsPush remove tags to a pending image
func UpdateImageTagsPull(mongoClient *mongo.Client, body types.BodyUpdateImageTagsPull) (*int64, error) {
	if body.ID == primitive.NilObjectID {
		return nil, errors.New("body not valid, ID empty")
	}
	collectionImagesPending := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("IMAGES_PENDING_COLLECTION"))
	return mongodb.UpdateImageTagsPull(collectionImagesPending, body)
}
