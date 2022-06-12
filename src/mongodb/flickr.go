package mongodb

import (
	"scrapper/src/types"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"fmt"
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

func FindImageId(collection *mongo.Collection, flickrId string) (*types.Image, error) {

	var image types.Image
	query := bson.M{"flickr_id": flickrId}

	options := options.FindOne().
		SetProjection(bson.M{
			"_id": 1,
		})

	err := collection.FindOne(context.TODO(), query, options).Decode(&image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}
