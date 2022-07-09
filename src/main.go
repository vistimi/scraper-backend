package main

import (
	"github.com/gin-gonic/gin"

	"scraper/src/routes"
	"scraper/src/utils"

	"net/http"

	"scraper/src/mongodb"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-contrib/cors"
)

func main() {

	mongoClient := mongodb.Connect()

	router := gin.Default()
	router.Use(cors.Default())

	// routes for one image pending or wanted
	router.Static("/image/file", utils.DotEnvVariable("IMAGE_PATH"))
	router.GET("/image/:id/:collection", wrapperHandlerURI(mongoClient, routes.FindImage))
	router.PUT("/image/tags/push", wrapperHandlerBody(mongoClient, routes.UpdateImageTagsPush))
	router.PUT("/image/tags/pull", wrapperHandlerBody(mongoClient, routes.UpdateImageTagsPull))
	router.PUT("/image/crop", wrapperHandlerBody(mongoClient, mongodb.UpdateImageCrop))
	router.POST("/image/crop", wrapperHandlerBody(mongoClient, mongodb.CreateImageCrop))
	router.POST("/image/transfer", wrapperHandlerBody(mongoClient, mongodb.TransferImage))
	router.DELETE("/image/:id", wrapperHandlerURI(mongoClient, routes.RemoveImageAndFile))

	// routes for multiple images pending or wanted
	router.GET("/images/id/:origin/:collection", wrapperHandlerURI(mongoClient, routes.FindImagesIDs))

	// routes for one image unwanted
	router.POST("/image/unwanted", wrapperHandlerBody(mongoClient, mongodb.InsertImageUnwanted))
	router.DELETE("/image/unwanted/:id", wrapperHandlerURI(mongoClient, routes.RemoveImage))

	// routes for multiple images unwanted
	router.GET("/images/unwanted", wrapperHandler(mongoClient, routes.FindImagesUnwanted))

	// routes for one tag
	router.POST("/tag/wanted", wrapperHandlerBody(mongoClient, mongodb.InsertTagWanted))
	router.POST("/tag/unwanted", wrapperHandlerBody(mongoClient, mongodb.InsertTagUnwanted))
	router.DELETE("/tag/wanted/:id", wrapperHandlerURI(mongoClient, routes.RemoveTagWanted))
	router.DELETE("/tag/unwanted/:id", wrapperHandlerURI(mongoClient, routes.RemoveTagUnwanted))

	// routes for multiple tags
	router.GET("/tags/wanted", wrapperHandler(mongoClient, mongodb.TagsWanted))
	router.GET("/tags/unwanted", wrapperHandler(mongoClient, mongodb.TagsUnwanted))

	// routes for one user unwanted
	router.POST("/user/unwanted", wrapperHandlerBody(mongoClient, mongodb.InsertUserUnwanted))
	router.DELETE("/user/unwanted/:id", wrapperHandlerURI(mongoClient, routes.RemoveUserUnwanted))

	// routes for multiplt users unwanted
	router.GET("/users/unwanted", wrapperHandler(mongoClient, mongodb.UsersUnwanted))

	// routes for scraping the internet
	router.POST("/search/flickr/:quality", wrapperHandlerURI(mongoClient, routes.SearchPhotosFlickr))
	router.POST("/search/unsplash/:quality", wrapperHandlerURI(mongoClient, routes.SearchPhotosUnsplash))
	router.POST("/search/pexels/:quality", wrapperHandlerURI(mongoClient, routes.SearchPhotosPexels))

	// start the backend
	router.Run("localhost:8080")
}

type mongoSchema interface {
	*mongo.Client
}

// wrapper for the response with argument
func wrapperResponseArg[M mongoSchema, A any, R any](c *gin.Context, f func(mongo M, arg A) (R, error), mongo M, arg A) {
	res, err := f(mongo, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// wrapper for the response
func wrapperResponse[M mongoSchema, R any](c *gin.Context, f func(mongo M) (R, error), mongo M) {
	res, err := f(mongo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// wrapper for the ginHandler with body with collectionName
func wrapperHandlerBody[B any, R any](mongoClient *mongo.Client, f func(mongo *mongo.Client, body B) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body B
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperResponseArg(c, f, mongoClient, body)
	}
}

// wrapper for the ginHandler with URI
func wrapperHandlerURI[P any, R any](mongoClient *mongo.Client, f func(mongo *mongo.Client, params P) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperResponseArg(c, f, mongoClient, params)
	}
}

// wrapper for the ginHandler
func wrapperHandler[R any](mongoClient *mongo.Client, f func(mongo *mongo.Client) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		wrapperResponse(c, f, mongoClient)
	}
}
