package main

import (
	"fmt"
	"scraper/src/mongodb"
	"scraper/src/router"
	"scraper/src/utils"
)

func main() {
	fmt.Println("Starting the server")
	mongoClient := mongodb.ConnectMongoDB()
	s3Client := utils.ConnectS3()
	_ = router.Router(mongoClient, s3Client)
}
