package mongodb

import (
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/mongo"

	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() *mongo.Client {

	uri := utils.DotEnvVariable("MONGODB_URI")

	// ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	// defer cancel()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return client

}
