package resource

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/Aman-Shitta/slotbook/internal/responses"
	"github.com/gin-gonic/gin"
)

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

func (rl *ResourcesListResponse) FilterByKind(name string) {

	var FilteredResources []ResourceResponse
	for _, item := range rl.Resources {
		if item.Kind == name {
			FilteredResources = append(FilteredResources, item)
		}
	}
	rl.Resources = FilteredResources
}

type ResourceUpdateRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty"`
	Kind     *string `json:"kind" binding:"omitempty,oneof=room studio desk"`
	Capacity *int    `json:"capacity,omitempty" binding:"omitempty,lte=500,gt=0"`
}

var ResourcesData = make(map[string]Resource)

func ResourcesHandler(c *gin.Context) {

	var response responses.Response

	switch c.Request.Method {
	case "POST":
		var req ResourceCreateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error = fmt.Sprintf("Something went wrong : %s", err.Error())
			c.JSON(400, response)
			return
		}

		resourceC, err := CreateNew(req.Name, req.Kind, req.Capacity)

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

		kind := c.Query("kind")

		for _, item := range ResourcesData {

			data.Resources = append(
				data.Resources,
				ResourceResponse{
					Id:        item.Id,
					Name:      item.Name,
					Kind:      item.Kind,
					Capacity:  item.Capacity,
					CreatedAt: item.Created_at,
				},
			)
		}

		if kind != "" {
			data.FilterByKind(kind)
		}

		slices.SortFunc(data.Resources, func(A, B ResourceResponse) int { return cmp.Compare(A.Name, B.Name) })

		response.Success = true
		response.Data = data
		response.Meta = &responses.Meta{Total: len(data.Resources), TotalPages: min(1, len(data.Resources))}

		c.JSON(200, response)
		return
	}
}

func ResourceHandler(c *gin.Context) {

	var response responses.Response

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
		// validate for explicit null vals
		var raw map[string]json.RawMessage

		if err := c.ShouldBindBodyWithJSON(&raw); err != nil {
			response.Error = fmt.Sprintf("Something went wrong : %s", err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}

		for _, key := range []string{"name", "kind", "capacity"} {
			if v, ok := raw[key]; ok && string(v) == "null" {
				response.Error = fmt.Sprintf("%s cannot be null", key)
				c.JSON(http.StatusUnprocessableEntity, response)
				return
			}
		}

		var resourceUpdate ResourceUpdateRequest
		if err := c.ShouldBindBodyWithJSON(&resourceUpdate); err != nil {
			response.Error = fmt.Sprintf("Something went wrong : %s", err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}

		fmt.Println("[+] resourceUpdate [+] ", resourceUpdate)
		if resource, ok := ResourcesData[id]; ok {

			if resourceUpdate.Name != nil && *resourceUpdate.Name != "" {
				resource.Name = *resourceUpdate.Name
			}

			if resourceUpdate.Kind != nil && *resourceUpdate.Kind != "" {
				resource.Kind = *resourceUpdate.Kind
			}

			if resourceUpdate.Capacity != nil && *resourceUpdate.Capacity != 0 {
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
