package server

import (
	routerGin "scraper-backend/src/driver/server/gin"
	"scraper-backend/src/util"

	"github.com/gin-gonic/gin"
)

func Constructor(cfg util.Config) *gin.Engine {
	return routerGin.Router(cfg)
}
