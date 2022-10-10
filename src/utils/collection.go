package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func ImagesCollection(mongoClient *mongo.Client, collection string) (*mongo.Collection, error) {
	switch collection {
	case "production":
		return mongoClient.Database(GetEnvVariable("SCRAPER_DB")).Collection(GetEnvVariable("PRODUCTION")), nil
	case "validation":
		return mongoClient.Database(GetEnvVariable("SCRAPER_DB")).Collection(GetEnvVariable("VALIDATION")), nil
	case "pending":
		return mongoClient.Database(GetEnvVariable("SCRAPER_DB")).Collection(GetEnvVariable("PENDING")), nil
	case "undesired":
		return mongoClient.Database(GetEnvVariable("SCRAPER_DB")).Collection(GetEnvVariable("UNDESIRED")), nil
	default:
		return nil, fmt.Errorf("`%s` does not exist for selecting the images collection. Choose `%s`, `%s` or `%s`",
			collection, "wanted", "pending", "unwanted")
	}
}
