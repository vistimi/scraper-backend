package adapter

import "github.com/gin-gonic/gin"

type DriverServerGin interface {
	Router(port int, healthCheckPath string) *gin.Engine
}
