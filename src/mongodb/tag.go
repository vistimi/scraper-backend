package mongodb

import (
	"errors"
	"fmt"
	"scrapper/src/types"
	"scrapper/src/utils"
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

	now := time.Now()
	document.CreationDate = &now
	res, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

type ReturnInsertTagUnwanted struct {
	InsertedTagId interface{}
	DeletedImageCount int64
}
func InsertTagUnwanted(mongoClient *mongo.Client, document types.Tag) (interface{}, error) {
	collectionTagsUnwated := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	insertedId, err := InsertTag(collectionTagsUnwated, document)
	if err != nil {
		return nil, err
	}

	query:= bson.M{
		"tags.name": document.Name,
	}
	collections := utils.ImageCollections(mongoClient)
	var deletedCount int64 = 0
	for _, collection := range collections {
		res, err := collection.DeleteMany(context.TODO(), query)
		if err != nil {
			return nil, err
		}
		deletedCount += res.DeletedCount
    }
	
	ids := ReturnInsertTagUnwanted{
		InsertedTagId: insertedId,
		DeletedImageCount: deletedCount,
	}
	return ids, nil
}

func InsertTagWanted(mongoClient *mongo.Client, document types.Tag) (interface{}, error) {
	collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	return InsertTag(collection, document)
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

func TagsWanted (mongoClient *mongo.Client) ([]types.Tag, error) {
	collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	return FindTags(collection)
}

func TagsUnwanted (mongoClient *mongo.Client) ([]types.Tag, error) {
	collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	return FindTags(collection)
}

func FindTags(collection *mongo.Collection) ([]types.Tag, error) {
	query := bson.D{}
	cursor, err := collection.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tags []types.Tag
	if err = cursor.All(context.TODO(), &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func FindTagName(collection *mongo.Collection, tagName string) (*types.Tag, error) {
	var tag types.Tag
	query := bson.M{"name": tagName}
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
