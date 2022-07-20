package main

import (
	"scraper/src/mongodb"
	"scraper/src/router"
)

func main() {
	mongoClient := mongodb.ConnectMongoDB()
	_ = router.Router(mongoClient)
}
