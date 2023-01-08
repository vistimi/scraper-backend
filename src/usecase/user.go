package usecase

import (
	"scraper-backend/src/mongodb"
	"scraper-backend/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	interfaceEntity "scraper-backend/src/entity/interface"
)

type usecaseUser struct {
}

func (u *usecaseUser) Contructor() interfaceEntity.UsecaseUser {
	return &usecaseUser{}
}

func RemoveUserUnwanted(mongoClient *mongo.Client, params ParamsRemoveTag) (*int64, error) {
	collectionUsersUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("USERS_UNDESIRED_COLLECTION"))
	userID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveUser(collectionUsersUnwanted, userID)
}
