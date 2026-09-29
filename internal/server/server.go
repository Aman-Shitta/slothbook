package server

import (
	"github.com/Aman-Shitta/slotbook/internal/booking"
	"github.com/Aman-Shitta/slotbook/internal/resource"
	"github.com/Aman-Shitta/slotbook/internal/responses"
	"github.com/gin-gonic/gin"
)

var V1_APPVERSION string = "v0.1.0"

func SetupRouter(port string) *gin.Engine {

	servEngine := gin.Default()

	// V1 URL grouping
	{
		v1 := servEngine.Group("/v1")

		// System
		v1.GET("/healthz", func(c *gin.Context) {
			c.JSON(
				200,
				responses.Response{
					Success: true,
					Data: gin.H{
						"status":  "healthy",
						"version": V1_APPVERSION,
					},
				},
			)
		})

		// Resources
		v1.Match([]string{"GET", "POST"}, "/resources/", resource.ResourcesHandler)
		v1.Match([]string{"GET", "DELETE", "PATCH"}, "/resource/:id", resource.ResourceHandler)

		// Bookings
		v1.Match([]string{"GET", "POST"}, "/bookings/", booking.BookingsHandler)
		v1.GET("/resources/:resourceId/bookings/", booking.ResourceBookingsHandler)
		v1.POST("/bookings/:id/cancel/", booking.CancelBookingHandler)
	}

	return servEngine
}
