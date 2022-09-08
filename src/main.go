package main

import (
	"log"
	"scraper/src/mongodb"
	"scraper/src/router"
	"scraper/src/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {
	var cfg aws.Config
	switch utils.DotEnvVariable("ENV") {
	case "production":
		cfg = utils.AwsS3()
	case "staging":
		log.Fatal("ENV staging not implemented yet")
	case "development":
		log.Fatal("ENV development not implemented yet")
	case "local":
		utils.LoadEnvVariables(".env")
		cfg = utils.LocalS3()
	default:
		log.Fatal("ENV variable is either production, staging, development or local")
	}
	s3Client := utils.ConnectS3(cfg)
	mongoClient := mongodb.ConnectMongoDB()
	_ = router.Router(mongoClient, s3Client)
}
