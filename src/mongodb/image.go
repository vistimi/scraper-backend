package mongodb

import (
	"errors"

	"scrapper/src/types"

	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"fmt"

	"time"

	"path/filepath"

	"os"
)

// InsertImage insert an image in its collection
func InsertImage(collection *mongo.Collection, document types.Image) (primitive.ObjectID, error) {
	res, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("Safecast of ObjectID did not work")
	}
	return insertedID, nil
}

// RemoveImage remove an image based on its mongodb id
func RemoveImage(collection *mongo.Collection, id primitive.ObjectID, origin string) (*int64, error) {
	query := bson.M{"_id": id, "origin": origin}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
}

// RemoveImageAndFile remove an image based on its mongodb id and remove its file
func RemoveImageAndFile(collection *mongo.Collection, origin string, id primitive.ObjectID) (*int64, error) {
	image, err := FindOne[types.Image](collection, bson.M{"_id": id, "origin": origin})
	if err != nil {
		return nil, fmt.Errorf("FindImageByID has failed: %v", err)
	}
	deletedCount, err := RemoveImage(collection, id, origin)
	if err != nil {
		return nil, fmt.Errorf("RemoveImage has failed: %v", err)
	}
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := fmt.Sprintf(filepath.Join(folderDir, origin, image.Path))
	err = os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("os.Remove has failed: %v", err)
	}
	return deletedCount, nil
}

func RemoveImagesAndFilesOneOrigin(mongoClient *mongo.Client, origin string, query bson.M, options *options.FindOptions) (*int64, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("IMAGES_COLLECTION"))
	var deletedCount int64
	images, err := FindMany[types.Image](collectionImages, query, options)
	if err != nil {
		return nil, fmt.Errorf("FindImagesIDs has failed: %v", err)
	}
	for _, image := range images {
		deletedOne, err := RemoveImageAndFile(collectionImages, origin, image.ID)
		if err != nil {
			return nil, fmt.Errorf("RemoveImageAndFile has failed: %v", err)
		}
		deletedCount += *deletedOne
	}
	return &deletedCount, nil
}

// Remove all the images in DB and their related file matching the query and options given, for all origins
func RemoveImagesAndFilesAllOrigins(mongoClient *mongo.Client, query bson.M, options *options.FindOptions) (*int64, error) {

	imageOrigins := utils.ImageOrigins()
	var deletedCount int64
	for _, origin := range imageOrigins {
		count, err := RemoveImagesAndFilesOneOrigin(mongoClient, origin, query, options)
		if err != nil {
			return nil, fmt.Errorf("RemoveImageAndFile has failed: %v", err)
		}
		deletedCount += *count
	}
	return &deletedCount, nil
}

// UpdateImage add tags to an image based on its mongodb id
func UpdateImage(collection *mongo.Collection, body types.BodyUpdateImage) (*types.Image, error) {
	query := bson.M{"_id": body.ID}
	if body.Tags != nil {
		for i := 0; i < len(body.Tags); i++ {
			tag := &body.Tags[i]
			now := time.Now()
			tag.CreationDate = &now
			fmt.Println(tag)
		}
		update := bson.M{
			"$push": bson.M{
				"tags": bson.M{"$each": body.Tags},
			},
		}
		_, err := collection.UpdateOne(context.TODO(), query, update)
		if err != nil {
			return nil, err
		}
	}
	return FindOne[types.Image](collection, bson.M{"_id": body.ID})
}
