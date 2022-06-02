package main

import (
	// "github.com/gin-gonic/gin"

	"dressme-scrapper/src/routes/flickr"

	// "net/http"

	// "dressme-scrapper/src/mongodb"
)

func main() {

	// collection := mongodb.Connect()

	// client := flickr.HttpClient()
	flickr.SearchPhoto("4", "model")
	// response := flickr.SendRequest(client, http.MethodPost)
	// fmt.Println(response)

	// router := gin.Default()
	// router.GET("/albums", routes.GetAlbums)
	// router.GET("/albums/:id", routes.GetAlbumByID)
	// router.POST("/albums", routes.PostAlbums)

	// router.Run("localhost:8080")
}
