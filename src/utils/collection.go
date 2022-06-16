package utils

import (
	"go.mongodb.org/mongo-driver/mongo"

	"errors"
)

func ImageCollectionSelection (mongoClient *mongo.Client, selection string) (*mongo.Collection, error) {
	switch selection {
	case "flickr":
		return mongoClient.Database(DotEnvVariable("SCRAPPER_DB")).Collection(DotEnvVariable("FLICKR_COLLECTION")), nil
	default:
		return nil, errors.New("Collection does not exist!")
	}
}

func ImageCollections (mongoClient *mongo.Client) map[string]*mongo.Collection {
	collections := make(map[string]*mongo.Collection)
	collections["flickr"] = mongoClient.Database(DotEnvVariable("SCRAPPER_DB")).Collection(DotEnvVariable("FLICKR_COLLECTION"))
	return collections
}