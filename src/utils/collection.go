package utils

import (
	"go.mongodb.org/mongo-driver/mongo"

	"errors"
)

// Return the MongoDB collection matching the desired selection
func ImageCollectionSelection (mongoClient *mongo.Client, selection string) (*mongo.Collection, error) {
	switch selection {
	case "flickr":
		return mongoClient.Database(DotEnvVariable("SCRAPPER_DB")).Collection(DotEnvVariable("FLICKR_COLLECTION")), nil
	default:
		return nil, errors.New("Collection does not exist")
	}
}

// Returns a map of all collections with images
func ImageCollections (mongoClient *mongo.Client) map[string]*mongo.Collection {
	collections := make(map[string]*mongo.Collection)
	collections["flickr"] = mongoClient.Database(DotEnvVariable("SCRAPPER_DB")).Collection(DotEnvVariable("FLICKR_COLLECTION"))
	return collections
}
