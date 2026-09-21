package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(port string) *gin.Engine {

	servEngine := gin.Default()

	servEngine.GET("/healthz", HealthHandler)

	servEngine.Match([]string{"GET", "POST"}, "/v1/resources/", ResourcesHandler)

	servEngine.Match([]string{"GET", "DELETE", "PATCH"}, "/v1/resource/:id", ResourceHandler)
	return servEngine
}
