package main

import (
	"log"
	"scraper-backend/src/driver/server"
	"scraper-backend/src/util"
)

func main() {
	config, err := util.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	_ = server.Constructor(mongoClient, s3Client)
}
