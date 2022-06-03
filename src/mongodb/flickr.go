package mongodb

import (
	"dressme-scrapper/src/types"

	"gopkg.in/mgo.v2"
	// "go.mongodb.org/mongo-driver/bson/primitives"
)

func InsertImage(collection *mgo.Collection, document types.FlickrImage) error {
	return collection.Insert(document)
}
