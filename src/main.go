package main

import (
	"github.com/gin-gonic/gin"

	"dressme-scrapper/src/routes"

	"net/http"

	"dressme-scrapper/src/mongodb"
)

func main() {

	collection := mongodb.Connect()

	client := routes.HttpClient()
	response := routes.SendRequest(client, http.MethodPost)

	// router := gin.Default()
	// router.GET("/albums", routes.GetAlbums)
	// router.GET("/albums/:id", routes.GetAlbumByID)
	// router.POST("/albums", routes.PostAlbums)

	// router.Run("localhost:8080")
}
