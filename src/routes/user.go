package routes

import (
	"scrapper/src/mongodb"
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RemoveUserUnwanted(mongoClient *mongo.Client, params ParamsRemoveTag) (*int64, error) {
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("USERS_UNWANTED_COLLECTION"))
	userID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveUser(collectionTagsUnwanted, userID)
}