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

func InsertImage(collection *mongo.Collection, document types.Image) (primitive.ObjectID, error) {
	res, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		message := fmt.Sprintf("Safecast of ObjectID did not work")
		return primitive.NilObjectID, errors.New(message)
	}
	return insertedId, nil
}

func FindImageByFLickrId(collection *mongo.Collection, flickrId string) (*types.Image, error) {
	var image types.Image
	query := bson.M{"flickrId": flickrId}
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

func FindImagesIds(collection *mongo.Collection, query bson.M) ([]types.Image, error) {
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

func FindImage(collection *mongo.Collection, id primitive.ObjectID) (*types.Image, error) {
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

func RemoveImage(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	query := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
}

func RemoveImageAndFile(collection *mongo.Collection, collectionDir string, id primitive.ObjectID) (*int64, error) {
	image, err := FindImage(collection, id)
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

func UpdateImage(collection *mongo.Collection, body types.BodyUpdateImage) (*types.Image, error) {
	query := bson.M{"_id": body.Id}
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
	return FindImage(collection, body.Id)
}
