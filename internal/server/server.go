package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(port string) *gin.Engine {

	servEngine := gin.Default()

	servEngine.GET("/healthz", HealthHandler)

	return servEngine
}
