package mongodb

import (
	"dressme-scrapper/src/types"

	"go.mongodb.org/mongo-driver/mongo"

	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertTag(collection *mongo.Collection, document types.Tag) (interface{}, error) {
	res, err:= collection.InsertOne(context.TODO(), document)
	if err != nil { 
		return nil, err 
	}
	return res.InsertedID, nil
}

func FindTagsUnwanted(collection *mongo.Collection) ([]types.Tag, error) {
	cursor, err:= collection.Find(context.TODO(), bson.D{})
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