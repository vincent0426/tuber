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

type TripPassenger struct {
	TripID        uuid.UUID
	PassengerID   uuid.UUID
	SourceID      uuid.UUID
	DestinationID uuid.UUID
	Status        string
	CreatedAt     time.Time
}
type NewTripPassenger struct {
	TripID        uuid.UUID
	PassengerID   uuid.UUID
	SourceID      uuid.UUID
	DestinationID uuid.UUID
	Status        string
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

// ID:                   dbTripView.ID,
//
//	DriverName:           dbTripView.DriverName,
//	DriverBrand:          dbTripView.DriverBrand,
//	DriverModel:          dbTripView.DriverModel,
//	DriverColor:          dbTripView.DriverColor,
//	DriverPlate:          dbTripView.DriverPlate,
//	SourceName:           dbTripView.SourceName,
//	SourcePlaceID:        dbTripView.SourcePlaceID,
//	SourceLatitude:       dbTripView.SourceLatitude,
//	SourceLongitude:      dbTripView.SourceLongitude,
//	DestinationName:      dbTripView.DestinationName,
//	DestinationPlaceID:   dbTripView.DestinationPlaceID,
//	DestinationLatitude:  dbTripView.DestinationLatitude,
//	DestinationLongitude: dbTripView.DestinationLongitude,
//	Status:               dbTripView.Status,
//	StartTime:            dbTripView.StartTime.In(time.Local),
//	CreatedAt:            dbTripView.CreatedAt.In(time.Local),
type TripView struct {
	ID                   uuid.UUID
	DriverName           string
	DriverBrand          string
	DriverModel          string
	DriverColor          string
	DriverPlate          string
	SourceName           string
	SourcePlaceID        string
	SourceLatitude       float64
	SourceLongitude      float64
	DestinationName      string
	DestinationPlaceID   string
	DestinationLatitude  float64
	DestinationLongitude float64
	Status               string
	StartTime            time.Time
	CreatedAt            time.Time
}
