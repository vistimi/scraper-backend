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

	"golang.org/x/exp/slices"

	"regexp"
)

// InsertTag inserts unique tag, not matching clsoe ones from its collection and the other one
func InsertTag(thisCollection *mongo.Collection, otherCollection *mongo.Collection, body types.Tag) (interface{}, error) {
	// only add unique tag from this collection
	thisTags, err := FindTags(thisCollection)
	if err != nil {
		return nil, err
	}
	idx := slices.IndexFunc(thisTags, func(thisTag types.Tag) bool {
		regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, strings.ToLower(thisTag.Name))
		matched, err := regexp.Match(regexpMatch, []byte(body.Name))	// e.g. match if thisTag has `model` and bodyTag `models`
		if err != nil {
			return false
		}
		return matched
	})
	if idx != -1 {
		return nil, fmt.Errorf("your tag `%s` and the db tag `%s` are too closely related", body.Name, thisTags[idx].Name)
	}

	// only unique tag from the other collection
	otherTags, err := FindTags(otherCollection)
	if err != nil {
		return nil, err
	}
	idx = slices.IndexFunc(otherTags, func(otherTag types.Tag) bool {
		regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, strings.ToLower(otherTag.Name))
		matched, err := regexp.Match(regexpMatch, []byte(body.Name))	// e.g. match if otherTag has `model` and bodyTag `models`
		if err != nil {
			return false
		}
		return matched
	})
	if idx != -1 {
		return nil, fmt.Errorf("your tag `%s` and the db tag `%s` are too closely related", body.Name, thisTags[idx].Name)
	}

	// insert tag
	now := time.Now()
	body.CreationDate = &now
	body.Name = strings.ToLower(body.Name)
	res, err := thisCollection.InsertOne(context.TODO(), body)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// ReturnInsertTagUnwanted indicates how many images with the new unwanted tag have been removed
type ReturnInsertTagUnwanted struct {
	InsertedTagId     interface{}
	DeletedImageCount int64
}

// InsertTagUnwanted inserts the new unwanted tag and remove the images with it as well as the files
func InsertTagUnwanted(mongoClient *mongo.Client, body types.Tag) (interface{}, error) {
	if body.Name == "" || body.Origin == "" {
		return nil, errors.New("Some fields are empty!")
	}
	body.Name = strings.ToLower(body.Name)

	// insert the unwanted tag
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	collectionTagsWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	insertedId, err := InsertTag(collectionTagsUnwanted, collectionTagsWanted, body)
	if err != nil {
		return nil, err
	}

	// remove the images with that unwanted tag
	query := bson.M{
		"tags.name": body.Name,
	}
	imageCollections := utils.ImageCollections(mongoClient)
	var deletedCount int64
	for collectionName, collection := range imageCollections {
		images, err := FindImagesIds(collection, query)
		if err != nil {
			return nil, err
		}
		for _, image := range images {
			deletedOne, err := RemoveImageAndFile(collection, collectionName, image.ID)
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

// InsertTagWanted insert a new tag 
func InsertTagWanted(mongoClient *mongo.Client, document types.Tag) (interface{}, error) {
	collectionTagsWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	return InsertTag(collectionTagsWanted, collectionTagsUnwanted, document)
}

// TagsWanted find all the wanted tags
func TagsWanted(mongoClient *mongo.Client) ([]types.Tag, error) {
	collectionTagsWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	return FindTags(collectionTagsWanted)
}

// TagsUnwanted find all the wanted tags
func TagsUnwanted(mongoClient *mongo.Client) ([]types.Tag, error) {
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	return FindTags(collectionTagsUnwanted)
}

// RemoveTag remove a tag from its collection
func RemoveTag(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	query := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
}

// FindTags find all the tags in its collection
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

// TagsWanted find all the names of wanted tags
func TagsWantedNames(mongoClient *mongo.Client) ([]string, error) {
	collectionWantedTags := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
	res, err := FindTags(collectionWantedTags)
	if err != nil {
		return nil, fmt.Errorf("FindTags Wanted has failed: \n%v", err)
	}
	var wantedTags []string
	for _, tag := range res {
		wantedTags = append(wantedTags, strings.ToLower(tag.Name))
	}
	return wantedTags, nil
}

// TagsUnwantednames find all the names of wanted tags
func TagsUnwantedNames(mongoClient *mongo.Client) ([]string, error) {
	collectionUnwantedTags := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
	res, err := FindTags(collectionUnwantedTags)
	if err != nil {
		return nil, fmt.Errorf("FindTags Unwated has failed: \n%v", err)
	}
	var unwantedTags []string
	for _, tag := range res {
		unwantedTags = append(unwantedTags, strings.ToLower(tag.Name))
	}
	return unwantedTags, nil
}