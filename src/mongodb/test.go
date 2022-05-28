package mongodb

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test(collection *mongo.Collection) {

	user := bson.D{{"fullName", "User 1"}, {"age", 30}}
	inserted, err := collection.InsertOne(context.TODO(), user)
	// check for errors in the insertion
	if err != nil {
		panic(err)
	}
	fmt.Println(inserted.InsertedID)

	title := "Back to the Future"
	var found bson.M
	err = collection.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&found)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(found, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"age", bson.D{{"$gt", 25}}},
				},
			},
		},
	}

	update := bson.D{
		{"$set",
			bson.D{
				{"age", 40},
			},
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Println("Number of documents updated:", result.ModifiedCount)
}
