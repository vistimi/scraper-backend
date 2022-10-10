package router

import (
	"scraper/src/mongodb"
	"scraper/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ParamsRemoveTag struct {
	ID string `uri:"id" binding:"required"`
}

func RemoveTagWanted(mongoClient *mongo.Client, params ParamsRemoveTag) (*int64, error) {
	collectionTagsWanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("TAGS_DESIRED_COLLECTION"))
	tagID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveTag(collectionTagsWanted, tagID)
}

func RemoveTagUnwanted(mongoClient *mongo.Client, params ParamsRemoveTag) (*int64, error) {
	collectionTagsUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("TAGS_UNDESIRED_COLLECTION"))
	tagID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveTag(collectionTagsUnwanted, tagID)
}
