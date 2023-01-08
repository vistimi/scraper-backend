package router

import (
	"fmt"
	"net/http"
	"scraper-backend/src/mongodb"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(mongoClient *mongo.Client, s3Client *s3.Client) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// health check
	router.Any("/", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })

	// routes for one image pending or wanted
	// router.Static("/image/file", utils.DotEnvVariable("IMAGE_PATH"))

	router.GET("/image/file/:origin/:name/:extension", wrapperDataHandlerURIS3(s3Client, FindImageFile))
	router.GET("/image/:id/:collection", wrapperJSONHandlerURI(mongoClient, FindImage))
	router.PUT("/image/tags/push", wrapperJSONHandlerBody(mongoClient, UpdateImageTagsPush))
	router.PUT("/image/tags/pull", wrapperJSONHandlerBody(mongoClient, UpdateImageTagsPull))
	router.PUT("/image/crop", wrapperJSONHandlerBodyS3(s3Client, mongoClient, mongodb.UpdateImageCrop))
	router.POST("/image/crop", wrapperJSONHandlerBodyS3(s3Client, mongoClient, mongodb.CreateImageCrop))
	router.POST("/image/copy", wrapperJSONHandlerBodyS3(s3Client, mongoClient, mongodb.CopyImage))
	router.POST("/image/transfer", wrapperJSONHandlerBody(mongoClient, mongodb.TransferImage))
	router.DELETE("/image/:id", wrapperJSONHandlerURIS3(s3Client, mongoClient, RemoveImageAndFile))

	// routes for multiple images pending or wanted
	router.GET("/images/id/:origin/:collection", wrapperJSONHandlerURI(mongoClient, FindImagesIDs))

	// routes for one image unwanted
	router.POST("/image/unwanted", wrapperJSONHandlerBody(mongoClient, mongodb.InsertImageUnwanted))
	router.DELETE("/image/unwanted/:id", wrapperJSONHandlerURI(mongoClient, RemoveImage))

	// routes for multiple images unwanted
	router.GET("/images/unwanted", wrapperJSONHandler(mongoClient, FindImagesUnwanted))

	// routes for one tag
	router.POST("/tag/wanted", wrapperJSONHandlerBody(mongoClient, mongodb.InsertTagWanted))
	router.POST("/tag/unwanted", wrapperJSONHandlerBodyS3(s3Client, mongoClient, mongodb.InsertTagUnwanted))
	router.DELETE("/tag/wanted/:id", wrapperJSONHandlerURI(mongoClient, RemoveTagWanted))
	router.DELETE("/tag/unwanted/:id", wrapperJSONHandlerURI(mongoClient, RemoveTagUnwanted))

	// routes for multiple tags
	router.GET("/tags/wanted", wrapperJSONHandler(mongoClient, mongodb.TagsWanted))
	router.GET("/tags/unwanted", wrapperJSONHandler(mongoClient, mongodb.TagsUnwanted))

	// routes for one user unwanted
	router.POST("/user/unwanted", wrapperJSONHandlerBodyS3(s3Client, mongoClient, mongodb.InsertUserUnwanted))
	router.DELETE("/user/unwanted/:id", wrapperJSONHandlerURI(mongoClient, RemoveUserUnwanted))

	// routes for multiple users unwanted
	router.GET("/users/unwanted", wrapperJSONHandler(mongoClient, mongodb.UsersUnwanted))

	// routes for scraping the internet
	router.POST("/search/flickr/:quality", wrapperJSONHandlerURIS3(s3Client, mongoClient, SearchPhotosFlickr))
	router.POST("/search/unsplash/:quality", wrapperJSONHandlerURIS3(s3Client, mongoClient, SearchPhotosUnsplash))
	router.POST("/search/pexels/:quality", wrapperJSONHandlerURIS3(s3Client, mongoClient, SearchPhotosPexels))

	// start the backend
	router.Run("0.0.0.0:8080")
	return router
}

type DataSchema struct {
	dataType string
	dataFile []byte
}

type mongoSchema interface {
	*mongo.Client
}

// wrapper for the response with argument
func wrapperJSONResponseArg[M mongoSchema, A any, R any](c *gin.Context, f func(mongo M, arg A) (R, error), mongo M, arg A) {
	res, err := f(mongo, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func wrapperJSONResponseArgS3[M mongoSchema, A any, R any](c *gin.Context, f func(s3Client *s3.Client, mongo M, arg A) (R, error), s3Client *s3.Client, mongo M, arg A) {
	res, err := f(s3Client, mongo, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func wrapperDataResponseArgS3[A any](c *gin.Context, f func(s3Client *s3.Client, arg A) (*DataSchema, error), s3Client *s3.Client, arg A) {
	data, err := f(s3Client, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	switch data.dataType {
	case "jpg":
		c.Data(http.StatusOK, "image/jpeg", data.dataFile)
	case "jpeg":
		c.Data(http.StatusOK, "image/jpeg", data.dataFile)
	case "png":
		c.Data(http.StatusOK, "image/png", data.dataFile)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"status": fmt.Errorf("wrong content-type: %s", data.dataType)})
	}
}

// wrapper for the response
func wrapperJSONResponse[M mongoSchema, R any](c *gin.Context, f func(mongo M) (R, error), mongo M) {
	res, err := f(mongo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// wrapper for the ginHandler with body with collectionName
func wrapperJSONHandlerBody[B any, R any](mongoClient *mongo.Client, f func(mongo *mongo.Client, body B) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body B
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperJSONResponseArg(c, f, mongoClient, body)
	}
}

// wrapper for the ginHandler with URI
func wrapperJSONHandlerURI[P any, R any](mongoClient *mongo.Client, f func(mongo *mongo.Client, params P) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperJSONResponseArg(c, f, mongoClient, params)
	}
}

// wrapper for the ginHandler
func wrapperJSONHandler[R any](mongoClient *mongo.Client, f func(mongo *mongo.Client) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		wrapperJSONResponse(c, f, mongoClient)
	}
}

// wrapper for the ginHandler with body with collectionName
func wrapperJSONHandlerBodyS3[B any, R any](s3Client *s3.Client, mongoClient *mongo.Client, f func(s3Client *s3.Client, mongo *mongo.Client, body B) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body B
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperJSONResponseArgS3(c, f, s3Client, mongoClient, body)
	}
}

// wrapper for the ginHandler with URI
func wrapperJSONHandlerURIS3[P any, R any](s3Client *s3.Client, mongoClient *mongo.Client, f func(s3Client *s3.Client, mongo *mongo.Client, params P) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperJSONResponseArgS3(c, f, s3Client, mongoClient, params)
	}
}

func wrapperDataHandlerURIS3[P any](s3Client *s3.Client, f func(s3Client *s3.Client, params P) (*DataSchema, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperDataResponseArgS3(c, f, s3Client, params)
	}
}
