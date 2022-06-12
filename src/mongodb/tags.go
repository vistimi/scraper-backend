package mongodb

import (
	"errors"
	"fmt"
	"scrapper/src/types"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertTag(collection *mongo.Collection, document types.Tag) (interface{}, error) {
	tag, err := FindTagName(collection, document.Name)
	if err != nil {
		return nil, err
	}
	if tag != nil {
		return nil, errors.New(fmt.Sprintf("the tag `%s` already exists: %v", document.Name, tag))
	}

	document.CreationDate = time.Now()
	res, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// func InsertTags(collection *mongo.Collection, documents []types.Tag) ([]interface{}, error) {
// 	documentsInterfaces := make([]interface{}, len(documents))
// 	for i := range documents {
// 		documentsInterfaces[i] = documents[i]
// 	}
// 	res, err := collection.InsertMany(context.TODO(), documentsInterfaces)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res.InsertedIDs, nil
// }

func FindTags(collection *mongo.Collection) ([]types.Tag, error) {
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []types.Tag
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func FindTagName(collection *mongo.Collection, tagName string) (*types.Tag, error) {
	var tag types.Tag
	query := bson.M{"tagName": tagName}
	err := collection.FindOne(context.TODO(), query).Decode(&tag)
	switch err {
	case nil:
		return &tag, nil
	case mongo.ErrNoDocuments:
		return nil, nil
	default:
		return nil, err
	}
}
