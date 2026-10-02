package booking

import (
	"time"
)

type Rfc339Time struct {
	Created_at time.Time
	Updated_at time.Time
}

type Booking struct {
	Rfc339Time

	Id         string
	ResourceId string
	Title      string
	Start      time.Time
	End        time.Time
	Status     string
}
