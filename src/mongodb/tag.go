package mongodb

import (
	"errors"
	"fmt"
	"scrapper/src/types"
	"scrapper/src/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"context"

	"go.mongodb.org/mongo-driver/bson"

	"sort"
)

func InsertTag(collection *mongo.Collection, body types.Tag) (interface{}, error) {
	// sorted array of tags
	foundTags, err := FindTags(collection)
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, tag := range foundTags {
		tags = append(tags, strings.ToLower(tag.Name))
	}
	sort.Strings(tags)
	fmt.Print(tags)

	// check in same collection if one tag regexp match the new tag
	needles := []string{strings.ToLower(body.Name)}
	i, j, err := utils.ContainsRegExp(tags, needles)
	if err != nil {
		return nil, err
	}
	if i != nil && j != nil {
		return nil, errors.New(fmt.Sprintf("your tag `%s` and the db tag `%s` are too closely related", needles[*i], tags[*j]))
	}

	// insert tag
	now := time.Now()
	body.CreationDate = &now
	body.Name = strings.ToLower(body.Name)
	inserted, err := collection.InsertOne(context.TODO(), body)
	if err != nil {
		return nil, err
	}
	return inserted.InsertedID, nil
}

type ReturnInsertTagUnwanted struct {
	InsertedTagId     interface{}
	DeletedImageCount int64
}

func InsertTagUnwanted(mongoClient *mongo.Client, body types.Tag) (interface{}, error) {
	if body.Name == "" || body.Origin == "" {
		return nil, errors.New("Some fields are empty!")
	}
	body.Name = strings.ToLower(body.Name)

	collectionTagsUnwated := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	insertedId, err := InsertTag(collectionTagsUnwated, body)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"tags.name": body.Name,
	}
	imageCollections := utils.ImageCollections(mongoClient)
	var deletedCount int64 = 0
	for collectionName, collection := range imageCollections {
		images, err := FindImagesIds(collection, query)
		if err != nil {
			return nil, err
		}
		for _, image := range images {
			deletedOne, err := RemoveImageAndFile(collection, collectionName, image.Id)
			if err != nil {
				return nil, err
			}
			deletedCount += *deletedOne
		}
	}

	ids := ReturnInsertTagUnwanted{
		InsertedTagId:     insertedId,
		DeletedImageCount: deletedCount,
	}
	return ids, nil
}

func InsertTagWanted(mongoClient *mongo.Client, document types.Tag) (interface{}, error) {
	collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	return InsertTag(collection, document)
}

func TagsWanted(mongoClient *mongo.Client) ([]types.Tag, error) {
	collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	return FindTags(collection)
}

func TagsUnwanted(mongoClient *mongo.Client) ([]types.Tag, error) {
	collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	return FindTags(collection)
}

func RemoveTag(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	query := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
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
