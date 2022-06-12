package main

import (
	"github.com/gin-gonic/gin"

	"scrapper/src/routes/flickr"

	"net/http"

	"scrapper/src/mongodb"
	"scrapper/src/types"
	"scrapper/src/utils"
)

func main() {

	mongoClient := mongodb.Connect()

	router := gin.Default()

	type ParamsFlickr struct {
		Quality string `uri:"quality" binding:"required"`
	}
	router.POST("/search/flickr/:quality", func(c *gin.Context) {
		var params ParamsFlickr
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		insertedIds, err := flickr.SearchPhoto(params.Quality, mongoClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"flickr images insertedIds": insertedIds})
	})

	router.POST("/tag/wanted", func(c *gin.Context) {
		var body types.Tag
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("WANTED_TAGS_COLLECTION"))
		insertedId, err := mongodb.InsertTag(collection, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"unwated tag insertedId": insertedId})
	})

	router.POST("/tag/unwanted", func(c *gin.Context) {
		var body types.Tag
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		collection := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNWANTED_TAGS_COLLECTION"))
		insertedId, err := mongodb.InsertTag(collection, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"unwated tag insertedId": insertedId})
	})

	router.Run("localhost:8080")
}
