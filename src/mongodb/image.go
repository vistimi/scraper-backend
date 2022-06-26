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

// FindImageIDByOriginID an image mongodb id based on its originID
func FindImageIDByOriginID(collection *mongo.Collection, originID string) (*types.Image, error) {
	var image types.Image
	query := bson.M{"originID": originID}
	options := options.FindOne().
		SetProjection(bson.M{
			"_id": 1,
		})
	err := collection.FindOne(context.TODO(), query, options).Decode(&image)
	switch err {
	case nil:
		return &image, nil
	case mongo.ErrNoDocuments:
		return nil, nil
	default:
		return nil, err
	}
}

// FindImagesIDs find all images mongodb id
func FindImagesIDs(collection *mongo.Collection, query bson.M) ([]types.Image, error) {
	options := options.Find().
		SetProjection(bson.M{
			"_id": 1,
		})
	cursor, err := collection.Find(context.TODO(), query, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var images []types.Image
	if err = cursor.All(context.TODO(), &images); err != nil {
		return nil, err
	}
	return images, nil
}

// FindImageByID find an image by its mongodb id
func FindImageByID(collection *mongo.Collection, id primitive.ObjectID) (*types.Image, error) {
	query := bson.M{"_id": id}
	var image types.Image
	err := collection.FindOne(context.TODO(), query).Decode(&image)
	switch err {
	case nil:
		return &image, nil
	case mongo.ErrNoDocuments:
		return nil, nil
	default:
		return nil, err
	}
}

// RemoveImage remove an image based on its mongodb id
func RemoveImage(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	query := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
}

// RemoveImageAndFile remove an image based on its mongodb id and remove its file
func RemoveImageAndFile(collection *mongo.Collection, collectionDir string, id primitive.ObjectID) (*int64, error) {
	image, err := FindImageByID(collection, id)
	if err != nil {
		return nil, err
	}
	deletedCount, err := RemoveImage(collection, id)
	if err != nil {
		return nil, err
	}
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := fmt.Sprintf(filepath.Join(folderDir, collectionDir, image.Path))
	err = os.Remove(path)
	if err != nil {
		return nil, err
	}
	return deletedCount, nil
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
	return FindImageByID(collection, body.ID)
}
