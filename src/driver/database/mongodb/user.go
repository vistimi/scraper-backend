package mongodb

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"scraper-backend/src/types"
// 	"scraper-backend/src/utils"
// 	"strings"
// 	"time"

// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type ReturnInsertUserUnwanted struct {
// 	InsertedTagID     interface{}
// 	DeletedImageCount int64
// }

// // InsertUserUnwanted inserts the new unwanted user and remove the images with it as well as the files
// func InsertUserUnwanted(s3Client *s3.Client, mongoClient *mongo.Client, body types.User) (*ReturnInsertUserUnwanted, error) {
// 	if body.Name == "" || body.Origin == "" || body.OriginID == "" {
// 		return nil, errors.New("some fields are empty")
// 	}
// 	now := time.Now()
// 	body.CreationDate = &now
// 	body.Origin = strings.ToLower(body.Origin)

// 	// insert the unwanted user
// 	collectionUserUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("USERS_UNDESIRED_COLLECTION"))
// 	query := bson.M{"origin": body.Origin,
// 		"$or": bson.A{
// 			bson.M{"originID": body.OriginID},
// 			bson.M{"name": body.Name},
// 		},
// 	}
// 	res, err := collectionUserUnwanted.InsertOne(context.TODO(), body)
// 	if err != nil {
// 		return nil, fmt.Errorf("InsertOne has failed: %v", err)
// 	}

// 	// remove the images with that unwanted user
// 	query = bson.M{
// 		"user.origin": body.Origin,
// 		"$or": bson.A{
// 			bson.M{"user.originID": body.OriginID},
// 			bson.M{"user.name": body.Name},
// 		},
// 	}
// 	options := options.Find().SetProjection(bson.M{"_id": 1})
// 	deletedCount, err := RemoveImagesAndFiles(s3Client, mongoClient, query, options) // check in all origins
// 	if err != nil {
// 		return nil, fmt.Errorf("RemoveImagesAndFiles has failed: %v", err)
// 	}

// 	ids := ReturnInsertUserUnwanted{
// 		InsertedTagID:     res.InsertedID,
// 		DeletedImageCount: *deletedCount,
// 	}
// 	return &ids, nil
// }

// // RemoveUser remove a tag from its collection
// func RemoveUser(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
// 	query := bson.M{"_id": id}
// 	res, err := collection.DeleteOne(context.TODO(), query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res.DeletedCount, nil
// }

// // TagsUnwanted find all the wanted tags
// func UsersUnwanted(mongoClient *mongo.Client) ([]types.User, error) {
// 	collectionUsersUnwanted := mongoClient.Database(utils.GetEnvVariable("SCRAPER_DB")).Collection(utils.GetEnvVariable("USERS_UNDESIRED_COLLECTION"))
// 	return FindMany[types.User](collectionUsersUnwanted, bson.M{})
// }
