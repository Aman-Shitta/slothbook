package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Aman-Shitta/slotbook/internal/resource"
	"github.com/gin-gonic/gin"
)

var APPVERSION string = "v0.1.0"

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func HealthHandler(c *gin.Context) {
	c.JSON(
		200,
		Response{
			Success: true,
			Data: gin.H{
				"status":  "healthy",
				"version": APPVERSION,
			},
		},
	)
}

var ResourcesData = make(map[string]resource.Resource)

type ResourceCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Kind     string `json:"kind" binding:"required,oneof=room studio desk"`
	Capacity int    `json:"capacity" binding:"required"`
}

type ResourceResponse struct {
	Id string `json:"id,omitempty"`
	ResourceCreateRequest
	CreatedAt time.Time `json:"created_at"`
}

type ResourcesListResponse struct {
	Resources []ResourceResponse `json:"resources" binding:"required"`
}

type ResourceUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Kind     *string `json:"kind" binding:"omitempty,oneof=room studio desk"`
	Capacity *int    `json:"capacity,omitempty" binding:"lte=500,gt=0"`
}

func ResourcesHandler(c *gin.Context) {

	var response Response
	switch c.Request.Method {
	case "POST":
		var req ResourceCreateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error = fmt.Sprintf("Something went wrong : %s", err.Error())
			c.JSON(400, response)
			return
		}

		resourceC, err := resource.CreateNew(req.Name, req.Kind, req.Capacity)

		if err != nil {
			response.Error = fmt.Sprintf("Something went wrong: %s", err.Error())
			c.JSON(400, response)
			return
		}

		ResourcesData[resourceC.Id] = *resourceC
		response.Success = true
		response.Data = "Resource Created"

		c.JSON(201, response)
		return

	case "GET":

		var data ResourcesListResponse

		for _, resource := range ResourcesData {
			data.Resources = append(
				data.Resources,
				ResourceResponse{
					Id:        resource.Id,
					Name:      resource.Name,
					Kind:      resource.Kind,
					Capacity:  resource.Capacity,
					CreatedAt: resource.Created_at,
				},
			)
		}

		response.Success = true
		response.Data = data
		response.Meta = &Meta{Total: len(data.Resources), TotalPages: min(1, len(data.Resources))}

		c.JSON(200, response)
		return
	}
}

func ResourceHandler(c *gin.Context) {

	var response Response

	id := c.Param("id")

	switch c.Request.Method {
	case "GET":

		if resource, ok := ResourcesData[id]; ok {
			response.Success = true
			response.Data = ResourceResponse{
				Name:      resource.Name,
				Kind:      resource.Kind,
				Capacity:  resource.Capacity,
				CreatedAt: resource.Created_at,
			}

			c.JSON(200, response)
			return
		} else {
			response.Error = "Data not present"
			c.JSON(400, response)
			return
		}

	case "DELETE":
		delete(ResourcesData, id)
		c.JSON(http.StatusNoContent, nil)
		return

	case "PATCH":
		var resourceUpdate ResourceUpdateRequest
		if err := c.ShouldBindJSON(&resourceUpdate); err != nil {
			response.Error = fmt.Sprintf("Something went wrong : %s", err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if resource, ok := ResourcesData[id]; ok {

			if resourceUpdate.Name == nil {
				response.Error = "name cannot be null"
				c.JSON(http.StatusUnprocessableEntity, response)
				return

			} else if *resourceUpdate.Name != "" {
				resource.Name = *resourceUpdate.Name
			}

			if resourceUpdate.Kind == nil {
				response.Error = "kind cannot be null"
				c.JSON(http.StatusUnprocessableEntity, response)
				return

			} else if *resourceUpdate.Kind != "" {
				resource.Kind = *resourceUpdate.Kind
			}

			if resourceUpdate.Capacity == nil {
				response.Error = "capacity cannot be null"
				c.JSON(http.StatusUnprocessableEntity, response)
				return

			} else if *resourceUpdate.Capacity != 0 {
				resource.Capacity = *resourceUpdate.Capacity
			}

			ResourcesData[id] = resource

		} else {
			response.Error = "Object not found"
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response.Success = true
		response.Data = "Data Updated"
		c.JSON(http.StatusOK, response)
	}
}
