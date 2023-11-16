package location

import (
	"github.com/google/uuid"
)

type Location struct {
	ID      uuid.UUID
	Name    string
	PlaceID string
	Lat     float64
	Lon     float64
}

// NewLocation contains information needed to create a new location.
type NewLocation struct {
	Name    string
	PlaceID string
	Lat     float64
	Lon     float64
}
