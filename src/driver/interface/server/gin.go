package adapter

import "github.com/gin-gonic/gin"

type DriverServerGin interface {
	Router() *gin.Engine
}