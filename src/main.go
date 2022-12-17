package main

import (
	"log"
	"scraper/src/mongodb"
	"scraper/src/router"
	"scraper/src/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	var s3Client *s3.Client
	switch utils.GetEnvVariable("ENV") {
	case "production":
		s3Client = utils.AwsS3()
	case "local":
		utils.LoadEnvVariables("local.env")
		s3Client = utils.LocalS3()
	default:
		log.Fatal("ENV variable is not defined")
	}
	mongoClient := mongodb.ConnectMongoDB()
	_ = router.Router(mongoClient, s3Client)
}
