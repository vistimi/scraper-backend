package main

import (
	"fmt"
	"scraper/src/mongodb"
	"scraper/src/router"
)

func main() {
	fmt.Println("Starting the server")
	mongoClient := mongodb.ConnectMongoDB()
	_ = router.Router(mongoClient)
}
