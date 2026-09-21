package resource

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type rfc339 struct {
	Created_at time.Time
	Updated_at time.Time
}

type Resource struct {
	rfc339

	Id       string
	Name     string
	Kind     string
	Capacity int
}

func CreateNew(name string, kind string, capacity int) (*Resource, error) {

	// validate name to be 3-80 chars
	if len(name) < 3 || len(name) > 80 {
		return nil, errors.New("name should be 3-80 chars only")
	}

	// validate Kind to be room, desk or studio
	var validKind = []string{"room", "desk", "studio"}
	if !slices.Contains(validKind, strings.ToLower(kind)) {
		return nil, errors.New("kind, should be:  `room`, `desk`, `studio`")
	}

	// validate capacity to be 1-500
	if capacity > 500 || capacity < 1 {
		return nil, errors.New("capacity must be between 1-100")
	}

	newUUid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Resource{
		Id:         newUUid.String(),
		Name:       name,
		Kind:       kind,
		Capacity:   capacity,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}, nil

}
