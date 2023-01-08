package main

import (
	"log"
	"scraper-backend/src/mongodb"
	"scraper-backend/src/router"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// mongoClient := mongodb.ConnectMongoDB()
	_ = router.Router(mongoClient, s3Client)
}
