package mongodb

import (
	"context"
	"fmt"
	"scraper/src/types"
	"scraper/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type wrapperSchema interface {
	types.User | types.Tag | types.Image
}

// FindUser find a user based on either its originID or userName
func FindOne[T wrapperSchema](collection *mongo.Collection, query bson.M, options ...*options.FindOneOptions) (*T, error) {
	var found T
	var err error
	if len(options) != 0 {
		err = collection.FindOne(context.TODO(), query, options[0]).Decode(&found)
	} else {
		err = collection.FindOne(context.TODO(), query).Decode(&found)
	}
	switch err {
	case nil:
		return &found, nil
	case mongo.ErrNoDocuments:
		return nil, nil
	default:
		return nil, err
	}
}

// FindTags find all the tags in its collection
func FindMany[T wrapperSchema](collection *mongo.Collection, query bson.M, options ...*options.FindOptions) ([]T, error) {
	var cursor *mongo.Cursor
	var err error
	if len(options) != 0 {
		cursor, err = collection.Find(context.TODO(), query, options[0])
	} else {
		cursor, err = collection.Find(context.TODO(), query)
	}
	fmt.Printf("errrr: %T\n %v", err, utils.ToJSON(err))
	if err != nil{
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var found []T
	if err = cursor.All(context.TODO(), &found); err != nil {
		return nil, err
	}
	return found, nil
}

func InsertOne[T wrapperSchema](collection *mongo.Collection, body T, query bson.M) (interface{}, error) {
	// only add unique element in the collection
	found, err := FindOne[T](collection, query)
	if err != nil {
		return nil, err
	}
	if found != nil {
		return nil, fmt.Errorf(`The element %T exist already in the collection`, *new(T))
	}

	// insert element if unique
	res, err := collection.InsertOne(context.TODO(), body)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}