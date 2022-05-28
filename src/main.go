package main

import (
	"github.com/gin-gonic/gin"

	"dressme-scrapper/src/routes"

	"dressme-scrapper/src/mongodb"
)

func main() {

	collection := mongodb.Connect()

	router := gin.Default()
	router.GET("/albums", routes.GetAlbums)
	router.GET("/albums/:id", routes.GetAlbumByID)
	router.POST("/albums", routes.PostAlbums)

	router.Run("localhost:8080")
}
