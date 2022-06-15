package main

import (

	"github.com/gin-gonic/gin"

	"scrapper/src/routes/flickr"

	"net/http"

	"scrapper/src/mongodb"
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	mongoClient := mongodb.Connect()

	router := gin.Default()

	router.POST("/search/flickr/:quality", wrapperHandlerUri(mongoClient, flickr.SearchPhoto))

	router.POST("/tag/wanted", wrapperHandlerBody(mongoClient, "WANTED_TAGS_COLLECTION", mongodb.InsertTag))
	router.POST("/tag/unwanted", wrapperHandlerBody(mongoClient, "UNWANTED_TAGS_COLLECTION", mongodb.InsertTag))

	router.Run("localhost:8080")
}

type mongoSchema interface {
	*mongo.Client | *mongo.Collection
}

// wrapper for the response after the function run
func wrapperResponse[M mongoSchema, A any, R any](c *gin.Context, f func(mongo M, arg A) (R, error), mongo M, arg A) {
	res, err := f(mongo, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "res": res})
}

// wrapper for the ginHandler with body
func wrapperHandlerBody[B any, R any](mongoClient *mongo.Client, collectionName string, f func(mongo *mongo.Collection, body B) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body B
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable(collectionName))
		wrapperResponse(c, f, collection, body)
	}
}

// wrapper for the ginHandler with URI
func wrapperHandlerUri[P any, R any](mongoClient *mongo.Client, f func(mongo *mongo.Client, params P) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperResponse(c, f, mongoClient, params)
	}
}
