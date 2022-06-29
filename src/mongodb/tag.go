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
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"

	"go.mongodb.org/mongo-driver/bson"

	"golang.org/x/exp/slices"

	"regexp"
	"sort"
)

// InsertTag inserts unique tag, not matching clsoe ones from its collection and the other one
func InsertTag(thisCollection *mongo.Collection, otherCollection *mongo.Collection, body types.Tag) (interface{}, error) {
	// only add unique tag from this collection
	thisTags, err := FindMany[types.Tag](thisCollection, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("FindTags has failed: %v", err)
	}
	idx := slices.IndexFunc(thisTags, func(thisTag types.Tag) bool {
		regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, strings.ToLower(thisTag.Name))
		matched, err := regexp.Match(regexpMatch, []byte(body.Name)) // e.g. match if thisTag has `model` and bodyTag `models`
		if err != nil {
			return false
		}
		return matched
	})
	if idx != -1 {
		return nil, fmt.Errorf("your tag `%s` and the db tag `%s` are too closely related", body.Name, thisTags[idx].Name)
	}

	// only unique tag from the other collection
	otherTags, err := FindMany[types.Tag](otherCollection, bson.M{})
	if err != nil {
		return nil, err
	}
	idx = slices.IndexFunc(otherTags, func(otherTag types.Tag) bool {
		regexpMatch := fmt.Sprintf(`[\-\_\w\d]*%s[\-\_\w\d]*`, strings.ToLower(otherTag.Name))
		matched, err := regexp.Match(regexpMatch, []byte(body.Name)) // e.g. match if otherTag has `model` and bodyTag `models`
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
	body.Origin = strings.ToLower(body.Origin)
	res, err := thisCollection.InsertOne(context.TODO(), body)
	if err != nil {
		return nil, fmt.Errorf("InsertTag has failed: %v", err)
	}
	return res.InsertedID, nil
}

// ReturnInsertTagUnwanted indicates how many images with the new unwanted tag have been removed
type ReturnInsertTagUnwanted struct {
	InsertedTagID     interface{}
	DeletedImageCount int64
}

// InsertTagUnwanted inserts the new unwanted tag and remove the images with it as well as the files
func InsertTagUnwanted(mongoClient *mongo.Client, body types.Tag) (*ReturnInsertTagUnwanted, error) {
	if body.Name == "" || body.Origin == "" {
		return nil, errors.New("Some fields are empty!")
	}
	body.Name = strings.ToLower(body.Name)
	body.Origin = strings.ToLower(body.Origin)

	// insert the unwanted tag
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_UNWANTED_COLLECTION"))
	collectionTagsWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_WANTED_COLLECTION"))
	insertedID, err := InsertTag(collectionTagsUnwanted, collectionTagsWanted, body)
	if err != nil {
		return nil, fmt.Errorf("InsertTag has failed: %v", err)
	}

	// remove the images with that unwanted tag
	query := bson.M{"tags.name": body.Name}
	options := options.Find().SetProjection(bson.M{"_id": 1})
	deletedCount, err := RemoveImagesAndFilesAllOrigins(mongoClient, query, options)
	if err != nil {
		return nil, fmt.Errorf("RemoveImagesAndFiles has failed: %v", err)
	}

	ids := ReturnInsertTagUnwanted{
		InsertedTagID:     insertedID,
		DeletedImageCount: *deletedCount,
	}
	return &ids, nil
}

// InsertTagWanted insert a new tag
func InsertTagWanted(mongoClient *mongo.Client, document types.Tag) (interface{}, error) {
	collectionTagsWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_WANTED_COLLECTION"))
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_UNWANTED_COLLECTION"))
	return InsertTag(collectionTagsWanted, collectionTagsUnwanted, document)
}

// TagsWanted find all the wanted tags
func TagsWanted(mongoClient *mongo.Client) ([]types.Tag, error) {
	collectionTagsWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_WANTED_COLLECTION"))
	return FindMany[types.Tag](collectionTagsWanted, bson.M{})
}

// TagsUnwanted find all the wanted tags
func TagsUnwanted(mongoClient *mongo.Client) ([]types.Tag, error) {
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_UNWANTED_COLLECTION"))
	return FindMany[types.Tag](collectionTagsUnwanted, bson.M{})
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

// TagsWanted find all the names of wanted tags
func TagsWantedNames(mongoClient *mongo.Client) ([]string, error) {
	collectionTagsWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_WANTED_COLLECTION"))
	res, err := FindMany[types.Tag](collectionTagsWanted, bson.M{})
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
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_UNWANTED_COLLECTION"))
	res, err := FindMany[types.Tag](collectionTagsUnwanted, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("FindTags Unwated has failed: \n%v", err)
	}
	var unwantedTags []string
	for _, tag := range res {
		unwantedTags = append(unwantedTags, strings.ToLower(tag.Name))
	}
	return unwantedTags, nil
}

func TagsNames(mongoClient *mongo.Client) ([]string, []string, error) {
	unwantedTags, err := TagsUnwantedNames(mongoClient)
	if err != nil {
		return nil, nil, err
	}
	if (unwantedTags == nil) || (len(unwantedTags) == 0) {
		return nil, nil, errors.New("unwantedTags are empty")
	}
	sort.Strings(unwantedTags)

	wantedTags, err := TagsWantedNames(mongoClient)
	if err != nil {
		return nil, nil, err
	}
	if (wantedTags == nil) || (len(wantedTags) == 0) {
		return nil, nil, errors.New("wantedTags are empty")
	}
	sort.Strings(wantedTags)
	return unwantedTags, wantedTags, nil
}
