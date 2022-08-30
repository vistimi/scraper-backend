package main

import (
	"scraper/src/mongodb"
	"scraper/src/router"
	"scraper/src/utils"
)

func main() {
	mongoClient := mongodb.ConnectMongoDB()
	s3Client := utils.ConnectS3()
	_ = router.Router(mongoClient, s3Client)
}
