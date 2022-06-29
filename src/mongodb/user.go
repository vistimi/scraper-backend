package mongodb

import (
	"context"
	"errors"
	"fmt"
	"scrapper/src/types"
	"scrapper/src/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindUser find a user based on either its originID or userName
func FindUser(collection *mongo.Collection, origin string, originID string, name string) (*types.User, error) {
	var user types.User
	query := bson.M{
		"origin": origin,
		"$or": bson.A{
			bson.M{"originID": originID},
			bson.M{"name": name},
		},
	}
	err := collection.FindOne(context.TODO(), query).Decode(&user)
	switch err {
	case nil:
		return &user, nil
	case mongo.ErrNoDocuments:
		return nil, nil
	default:
		return nil, err
	}
}

// InsertUser inserts a unique user
func insertUser(userCollection *mongo.Collection, body types.User) (interface{}, error) {
	// only add unique user from this collection
	userFound, err := FindUser(userCollection, body.Origin, body.OriginID, body.Name)
	if err != nil {
		return nil, err
	}
	if userFound != nil {
		return nil, errors.New(`The user exist already in the collection`)
	}

	// insert user
	now := time.Now()
	body.CreationDate = &now
	body.Origin = strings.ToLower(body.Origin)
	res, err := userCollection.InsertOne(context.TODO(), body)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

type ReturnInsertUserUnwanted struct {
	InsertedTagID     interface{}
	DeletedImageCount int64
}

// InsertUserUnwanted inserts the new unwanted user and remove the images with it as well as the files
func InsertUserUnwanted(mongoClient *mongo.Client, body types.User) (*ReturnInsertUserUnwanted, error) {
	if body.Name == "" || body.Origin == "" || body.OriginID == "" {
		return nil, errors.New("Some fields are empty!")
	}
	body.Origin = strings.ToLower(body.Origin)

	// insert the unwanted user
	collectionUserUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("USERS_UNWANTED_COLLECTION"))
	insertedID, err := insertUser(collectionUserUnwanted, body)
	if err != nil {
		return nil, fmt.Errorf("insertUser has failed: %v", err)
	}

	// remove the images with that unwanted user
	query := bson.M{
		"user.origin": body.Origin,
		"$or": bson.A{
			bson.M{"user.originID": body.OriginID},
			bson.M{"user.name": body.Name},
		},
	}
	options := options.Find().SetProjection(bson.M{"_id": 1})
	deletedCount, err := RemoveImagesAndFilesOneOrigin(mongoClient, body.Origin, query, options)	// check in all origins
	if err != nil {
		return nil, fmt.Errorf("RemoveImagesAndFiles has failed: %v", err)
	}

	ids := ReturnInsertUserUnwanted{
		InsertedTagID:     insertedID,
		DeletedImageCount: *deletedCount,
	}
	return &ids, nil
}

// RemoveUser remove a tag from its collection
func RemoveUser(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	query := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
}

// TagsUnwanted find all the wanted tags
func UsersUnwanted(mongoClient *mongo.Client) ([]types.User, error) {
	collectionUsersUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("USERS_UNWANTED_COLLECTION"))
	return FindMany[types.User](collectionUsersUnwanted, bson.M{})
}