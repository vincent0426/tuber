package location

import (
	"github.com/google/uuid"
)

type QueryFilter struct {
	ID      *uuid.UUID `validate:"omitempty"`
	Name    *string    `validate:"omitempty"`
	PlaceID *string    `validate:"omitempty"`
	Lat     *float64   `validate:"omitempty"`
	Lon     *float64   `validate:"omitempty"`
}
