package gin

import (
	"context"
	"fmt"
	"net/http"

	interfaceAdapter "scraper-backend/src/adapter/interface"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type DriverServerGin struct {
	ControllerPicture interfaceAdapter.ControllerPicture
	ControllerTag     interfaceAdapter.ControllerTag
	ControllerUser    interfaceAdapter.ControllerUser
	ControllerFlickr interfaceAdapter.ControllerFlickr
	ControllerPexels interfaceAdapter.ControllerPexels
	ControllerUnsplash interfaceAdapter.ControllerUnsplash
}

// TODO: check Body and URI match path
func (d DriverServerGin) Router() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// health check
	router.Any("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })

	router.GET("/image/file/:origin/:name/:extension", wrapperDataHandlerURI(d.ReadPictureFile))
	router.GET("/image/:id/:collection", wrapperJSONHandlerURI(d.ReadPicture))
	router.PUT("/image/tag", wrapperJSONHandlerBody(d.UpdatePictureTag))    // TODO: changed URI
	router.DELETE("/image/tag", wrapperJSONHandlerBody(d.DeletePictureTag)) // TODO: changed URI
	router.PUT("/image/crop", wrapperJSONHandlerBody(d.UpdatePictureCrop))
	router.POST("/image/crop", wrapperJSONHandlerBody(d.CreatePictureCrop))
	router.POST("/image/copy", wrapperJSONHandlerBody(d.CreatePictureCopy))
	router.POST("/image/transfer", wrapperJSONHandlerBody(d.UpdatePictureTransfer))
	router.DELETE("/image/:id", wrapperJSONHandlerURI(d.DeletePictureAndFile))

	// routes for multiple images
	router.GET("/images/id/:origin/:collection", wrapperJSONHandlerURI(d.ReadPicturesID))

	// routes for one image unwanted
	router.POST("/image/unwanted", wrapperJSONHandlerBody(d.CreatePictureBlocked))
	router.DELETE("/image/unwanted", wrapperJSONHandlerURI(d.DeletePictureBlocked))

	// routes for multiple images unwanted
	router.GET("/images/unwanted", wrapperJSONHandler(d.ReadPicturesBlocked))

	// routes for one tag
	router.POST("/tag/wanted", wrapperJSONHandlerBody(d.CreateTag))
	router.POST("/tag/unwanted", wrapperJSONHandlerBody(d.CreateTagBlocked))
	router.DELETE("/tag/wanted/:id", wrapperJSONHandlerURI(d.DeleteTag))         
	router.DELETE("/tag/unwanted/:id", wrapperJSONHandlerURI(d.DeleteTagBlocked)) 

	// routes for multiple tags
	router.GET("/tags/wanted", wrapperJSONHandler(d.ReadTags))
	router.GET("/tags/unwanted", wrapperJSONHandler(d.ReadTagsBlocked))

	// routes for one user unwanted
	router.POST("/user/unwanted", wrapperJSONHandlerBody(d.CreateUserBlocked))
	router.DELETE("/user/unwanted/:id", wrapperJSONHandlerURI(d.DeleteUserBlocked))

	// routes for multiple users unwanted
	router.GET("/users/unwanted", wrapperJSONHandler(d.ReadUsers))

	// routes for scraping the internet
	router.POST("/search/flickr/:quality", wrapperJSONHandlerURI(d.SearchPhotosFlickr))
	router.POST("/search/unsplash/:quality", wrapperJSONHandlerURI(d.SearchPhotosUnsplash))
	router.POST("/search/pexels/:quality", wrapperJSONHandlerURI(d.SearchPhotosPexels))

	// start the backend
	router.Run("0.0.0.0:8080")
	return router
}

// JSON response

// Body
func wrapperJSONHandlerBody[B any, R any](f func(ctx context.Context, body B) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body B
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperJSONResponseArg(c, f, body)
	}
}

// URI
func wrapperJSONHandlerURI[P any, R any](f func(ctx context.Context, params P) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperJSONResponseArg(c, f, params)
	}
}

func wrapperJSONResponseArg[A any, R any](c *gin.Context, f func(ctx context.Context, arg A) (R, error), arg A) {
	res, err := f(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// No Body and URI
func wrapperJSONHandler[R any](f func(ctx context.Context) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		wrapperJSONResponse(c, f)
	}
}

func wrapperJSONResponse[R any](c *gin.Context, f func(ctx context.Context) (R, error)) {
	res, err := f(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// File response

type DataSchema struct {
	DataType string
	DataFile []byte
}

// URI
func wrapperDataHandlerURI[P any](f func(ctx context.Context, params P) (*DataSchema, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperDataResponseArg(c, f, params)
	}
}

func wrapperDataResponseArg[A any](c *gin.Context, f func(ctx context.Context, arg A) (*DataSchema, error), arg A) {
	data, err := f(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	switch data.DataType {
	case "jpg":
		c.Data(http.StatusOK, "image/jpeg", data.DataFile)
	case "jpeg":
		c.Data(http.StatusOK, "image/jpeg", data.DataFile)
	case "png":
		c.Data(http.StatusOK, "image/png", data.DataFile)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"status": fmt.Errorf("wrong content-type: %s", data.DataType)})
	}
}
