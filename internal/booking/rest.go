package booking

import (
	"fmt"
	"net/http"
	"time"
	"uuid"

	"github.com/Aman-Shitta/slotbook/internal/resource"
	responses "github.com/Aman-Shitta/slotbook/internal/responses"
	"github.com/gin-gonic/gin"
)

type BookingResponse struct {
	Id       string                    `json:"id"`
	Resource resource.ResourceResponse `json:"resource"`
	Title    string                    `json:"title"`
	Status   string                    `json:"status" binding:"oneof=confirmed cancelled"`
	Start    time.Time                 `json:"start"`
	End      time.Time                 `json:"end"`
}

type BookingResponseList struct {
	Bookings []BookingResponse `json:"bookings"`
}

type BookingCreateRequest struct {
	ResourceId string    `json:"resource_id" binding:"required"`
	Title      string    `json:"title" binding:"required"`
	Start      time.Time `json:"start" binding:"required"`
	End        time.Time `json:"end" binding:"required"`
}

var BookingsData = make(map[string]Booking)

// Checks weather a resource is already booked
// if booked returns list of conflicting booking resourceId
// else return null resp
func BookingConflictCheck(resourceId string, start, end time.Time) []string {
	var alreadyBooked []string

	for bookingId, booking := range BookingsData {
		if booking.ResourceId == resourceId {
			if (booking.Start.Compare(start) == -1 && booking.End.Compare(end) >= 0) || (booking.Start.Compare(end) == -1 && booking.End.Compare(end) >= 0) {
				alreadyBooked = append(alreadyBooked, bookingId)
			}
		}
	}
	fmt.Println("[+] ALready Booked [+]", alreadyBooked)
	return alreadyBooked
}

const (
	MAXBOOKINGDURATION  = 8 * 60 * 60 * time.Second
	MINNBOOKINGDURATION = 15 * 60 * time.Second
)

func BookingsHandler(c *gin.Context) {

	var response responses.Response

	switch c.Request.Method {
	case "GET":
		var BookingList BookingResponseList

		for _, booking := range BookingsData {
			bookedResource, ok := resource.ResourcesData[booking.ResourceId]
			fmt.Println("[+] Boked resource [+]", ok, booking.ResourceId)
			bookignResponse := BookingResponse{
				Id: booking.Id,
				Resource: resource.ResourceResponse{
					Id:        bookedResource.Id,
					Name:      bookedResource.Name,
					Kind:      bookedResource.Kind,
					Capacity:  bookedResource.Capacity,
					CreatedAt: bookedResource.Created_at,
				},
				Title:  booking.Title,
				Status: booking.Status,
				Start:  booking.Start,
				End:    booking.End,
			}

			BookingList.Bookings = append(
				BookingList.Bookings,
				bookignResponse,
			)

		}

		response.Success = true
		response.Data = BookingList
		response.Meta = &responses.Meta{
			Page:       1,
			Total:      len(BookingList.Bookings),
			TotalPages: min(1, len(BookingList.Bookings)),
		}

		c.JSON(http.StatusOK, response)
		return

	case "POST":
		var request BookingCreateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			response.Error = fmt.Sprintf("Invalid Data : %s", err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}
		var bookingResource *resource.Resource

		// validate Resource ID
		for _, availableResource := range resource.ResourcesData {
			if availableResource.Id == request.ResourceId {
				bookingResource = &availableResource
			}
		}

		if bookingResource == nil {
			response.Error = fmt.Sprintf("Resource `%s` not available to book", request.ResourceId)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		// validate start and end
		if time.Since(request.Start) >= 0 {
			response.Error = "Start Time cannot be in past"
			c.JSON(http.StatusBadGateway, response)
			return
		}

		bookingDuraction := request.End.Sub(request.Start)

		if bookingDuraction < MINNBOOKINGDURATION || bookingDuraction > MAXBOOKINGDURATION {
			response.Error = "The bookings must be min 15m and Max 8 hours"
			c.JSON(http.StatusBadGateway, response)
			return
		}

		//validate no two bookings for same resource in same time overlap
		if bookings := BookingConflictCheck(request.ResourceId, request.Start, request.End); len(bookings) != 0 {
			response.Error = "Resource not available for the given time"
			response.Data = map[string][]string{
				"bookings": bookings,
			}
			c.JSON(http.StatusConflict, response)
			return
		}
		// add booking
		bookingId := uuid.NewV7().String()
		BookingsData[bookingId] = Booking{
			Id:         bookingId,
			ResourceId: request.ResourceId,
			Title:      request.Title,
			Start:      request.Start,
			End:        request.End,
			Status:     "Confirmed",
			Created_at: time.Now(),
			Updated_at: time.Now(),
		}

		response.Success = true
		response.Data = fmt.Sprintf("Booking Created `%s`", bookingId)
		c.JSON(http.StatusOK, response)

	}
}

func ResourceBookingsHandler(c *gin.Context) {
	resourceId := c.Param("resourceId")

	var response responses.Response

	var bookinglist BookingResponseList

	for _, booking := range BookingsData {
		if booking.ResourceId == resourceId {
			bookingResource := resource.ResourcesData[resourceId]

			bookinglist.Bookings = append(
				bookinglist.Bookings,
				BookingResponse{
					Id: booking.Id,
					Resource: resource.ResourceResponse{
						Id:        bookingResource.Id,
						Name:      bookingResource.Name,
						Kind:      bookingResource.Kind,
						Capacity:  bookingResource.Capacity,
						CreatedAt: bookingResource.Created_at,
					},
					Title:  booking.Title,
					Status: booking.Status,
					Start:  booking.Start,
					End:    booking.End,
				},
			)
		}
		response.Success = true
		response.Data = bookinglist
		response.Meta = &responses.Meta{
			Total:      len(bookinglist.Bookings),
			TotalPages: min(1, len(bookinglist.Bookings)),
		}
		c.JSON(http.StatusOK, response)
	}
}

func CancelBookingHandler(c *gin.Context) {
	bookingId := c.Param("id")
	var response responses.Response

	booking, ok := BookingsData[bookingId]

	if !ok {
		response.Error = fmt.Sprintf("Booking `%s` do not exist.", bookingId)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// validate the booking end is in past cannont be cancelled
	if time.Now().Compare(booking.End) > 0 {
		response.Error = fmt.Sprintf("Booking `%s` cannot be cancelled", bookingId)
		c.JSON(http.StatusConflict, response)
		return
	}

	booking.Status = "cancelled"
	booking.Updated_at = time.Now()

	BookingsData[bookingId] = booking

	response.Success = true
	response.Data = fmt.Sprintf("Booking `%s` cancelled", bookingId)

	c.JSON(http.StatusOK, response)

}
