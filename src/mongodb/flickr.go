package mongodb

import (
	"errors"
	"scrapper/src/types"

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
