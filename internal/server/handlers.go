package server

import (
	"github.com/gin-gonic/gin"
)

var APPVERSION string = "v0.1.0"

func HealthHandler(c *gin.Context) {
	c.JSON(
		200,
		map[string]string{
			"status":  "healthy",
			"version": APPVERSION,
		},
	)
}
