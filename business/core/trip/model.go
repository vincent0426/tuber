package trip

import (
	"time"

	"github.com/google/uuid"
)

// User represents information about an individual user.
type Trip struct {
	ID             uuid.UUID
	DriverID       uuid.UUID
	PassengerLimit int
	SourceID       uuid.UUID
	DestinationID  uuid.UUID
	Status         string
	StartTime      time.Time
	CreatedAt      time.Time
}

// NewTrip contains information needed to create a new trip.
type NewTrip struct {
	DriverID       uuid.UUID
	PassengerLimit int
	SourceID       uuid.UUID
	DestinationID  uuid.UUID
	StartTime      time.Time
}

type UserTrip struct {
	TripID               uuid.UUID
	PassengerID          uuid.UUID
	StationSourceID      uuid.UUID
	StationDestinationID uuid.UUID
	PassengerStatus      string
	DriverID             uuid.UUID
	PassengerLimit       int
	SourceID             uuid.UUID
	DestinationID        uuid.UUID
	TripStatus           string
	StartTime            time.Time
	CreatedAt            time.Time
}
