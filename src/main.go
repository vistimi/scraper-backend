package main

import (
	"github.com/gin-gonic/gin"

	"scrapper/src/routes"

	"net/http"

	"scrapper/src/mongodb"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-contrib/cors"
)

func main() {

	mongoClient := mongodb.Connect()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/images/id/:origin", wrapperHandlerURI(mongoClient, routes.FindImagesIDs))

	router.Static("/image/file", "/home/olivier/dressme/images")
	router.GET("/image/:id", wrapperHandlerURI(mongoClient, routes.FindImage))
	router.PUT("/image", wrapperHandlerBody(mongoClient, routes.UpdateImage))
	router.DELETE("/image", wrapperHandlerBody(mongoClient, routes.RemoveImage))

	router.POST("/tag/wanted", wrapperHandlerBody(mongoClient, mongodb.InsertTagWanted))
	router.POST("/tag/unwanted", wrapperHandlerBody(mongoClient, mongodb.InsertTagUnwanted))
	router.DELETE("/tag/wanted/:id", wrapperHandlerURI(mongoClient, routes.RemoveTagWanted))
	router.DELETE("/tag/unwanted/:id", wrapperHandlerURI(mongoClient, routes.RemoveTagUnwanted))
	router.GET("/tags/wanted", wrapperHandler(mongoClient, mongodb.TagsWanted))
	router.GET("/tags/unwanted", wrapperHandler(mongoClient, mongodb.TagsUnwanted))

	router.POST("/search/flickr/:quality", wrapperHandlerURI(mongoClient, routes.SearchPhotosFlickr))
	router.POST("/search/unsplash", wrapperHandler(mongoClient, routes.SearchPhotosUnsplash))
	router.POST("/search/pexels", wrapperHandler(mongoClient, routes.SearchPhotosPexels))

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
