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

	"time"
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

func FindImagesIds(collection *mongo.Collection) ([]types.Image, error) {
	query := bson.D{}
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
