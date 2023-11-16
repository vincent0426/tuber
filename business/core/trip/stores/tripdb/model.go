package tripdb

import (
	"time"

	"github.com/TSMC-Uber/server/business/core/trip"
	"github.com/google/uuid"
)

// dbTrip represent the structure we need for moving data
// between the app and the database.
type dbTrip struct {
	ID             uuid.UUID `db:"id"`
	DriverID       uuid.UUID `db:"driver_id"`
	PassengerLimit int       `db:"passenger_limit"`
	SourceID       uuid.UUID `db:"source_id"`
	DestinationID  uuid.UUID `db:"destination_id"`
	Status         string    `db:"status"`
	StartTime      time.Time `db:"start_time"`
	CreatedAt      time.Time `db:"created_at"`
}

func toDBTrip(trip trip.Trip) dbTrip {
	return dbTrip{
		ID:             trip.ID,
		DriverID:       trip.DriverID,
		PassengerLimit: trip.PassengerLimit,
		SourceID:       trip.SourceID,
		DestinationID:  trip.DestinationID,
		StartTime:      trip.StartTime.UTC(),
		CreatedAt:      trip.CreatedAt.UTC(),
	}
}

func toCoreTrip(dbTrip dbTrip) trip.Trip {

	trip := trip.Trip{
		ID:             dbTrip.ID,
		DriverID:       dbTrip.DriverID,
		PassengerLimit: dbTrip.PassengerLimit,
		SourceID:       dbTrip.SourceID,
		DestinationID:  dbTrip.DestinationID,
		Status:         dbTrip.Status,
		StartTime:      dbTrip.StartTime.In(time.Local),
		CreatedAt:      dbTrip.CreatedAt.In(time.Local),
	}

	return trip
}

func toCoreTripSlice(dbTrips []dbTrip) []trip.Trip {
	trips := make([]trip.Trip, len(dbTrips))
	for i, dbTrip := range dbTrips {
		trips[i] = toCoreTrip(dbTrip)
	}
	return trips
}

// ------------------------------------------------------------
type dbUserTrip struct {
	TripID               uuid.UUID `db:"trip_id"`
	PassengerID          uuid.UUID `db:"passenger_id"`
	StationSourceID      uuid.UUID `db:"source_id"`
	StationDestinationID uuid.UUID `db:"destination_id"`
	PassengerStatus      string    `db:"tp_status"`
	DriverID             uuid.UUID `db:"driver_id"`
	PassengerLimit       int       `db:"passenger_limit"`
	SourceID             uuid.UUID `db:"source_id"`
	DestinationID        uuid.UUID `db:"destination_id"`
	TripStatus           string    `db:"trip_status"`
	StartTime            time.Time `db:"start_time"`
	CreatedAt            time.Time `db:"trip_created_at"`
}

func toCoreUserTrip(dbUserTrip dbUserTrip) trip.UserTrip {

	trip := trip.UserTrip{
		TripID:               dbUserTrip.TripID,
		PassengerID:          dbUserTrip.PassengerID,
		StationSourceID:      dbUserTrip.StationSourceID,
		StationDestinationID: dbUserTrip.StationDestinationID,
		PassengerStatus:      dbUserTrip.PassengerStatus,
		DriverID:             dbUserTrip.DriverID,
		PassengerLimit:       dbUserTrip.PassengerLimit,
		SourceID:             dbUserTrip.SourceID,
		DestinationID:        dbUserTrip.DestinationID,
		TripStatus:           dbUserTrip.TripStatus,
		StartTime:            dbUserTrip.StartTime.In(time.Local),
		CreatedAt:            dbUserTrip.CreatedAt.In(time.Local),
	}

	return trip
}

func toCoreUserTripSlice(dbUserTrips []dbUserTrip) []trip.UserTrip {
	trips := make([]trip.UserTrip, len(dbUserTrips))
	for i, dbTrip := range dbUserTrips {
		trips[i] = toCoreUserTrip(dbTrip)
	}
	return trips
}

// ------------------------------------------------------------
type dbTripPassenger struct {
	TripID        uuid.UUID `db:"trip_id"`
	PassengerID   uuid.UUID `db:"passenger_id"`
	SourceID      uuid.UUID `db:"source_id"`
	DestinationID uuid.UUID `db:"destination_id"`
	Status        string    `db:"tp_status"`
	CreatedAt     time.Time `db:"created_at"`
}

func toDBTripPassenger(tripPassenger trip.TripPassenger) dbTripPassenger {
	return dbTripPassenger{
		TripID:        tripPassenger.TripID,
		PassengerID:   tripPassenger.PassengerID,
		SourceID:      tripPassenger.SourceID,
		DestinationID: tripPassenger.DestinationID,
		Status:        tripPassenger.Status,
		CreatedAt:     tripPassenger.CreatedAt.UTC(),
	}
}

// ------------------------------------------------------------
type dbTripView struct {
	ID                   uuid.UUID `db:"id"`
	DriverName           string    `db:"driver_name"`
	DriverBrand          string    `db:"driver_brand"`
	DriverModel          string    `db:"driver_model"`
	DriverColor          string    `db:"driver_color"`
	DriverPlate          string    `db:"driver_plate"`
	SourceName           string    `db:"source_name"`
	SourcePlaceID        string    `db:"source_place_id"`
	SourceLatitude       float64   `db:"source_latitude"`
	SourceLongitude      float64   `db:"source_longitude"`
	DestinationName      string    `db:"destination_name"`
	DestinationPlaceID   string    `db:"destination_place_id"`
	DestinationLatitude  float64   `db:"destination_latitude"`
	DestinationLongitude float64   `db:"destination_longitude"`
	Status               string    `db:"status"`
	StartTime            time.Time `db:"start_time"`
	CreatedAt            time.Time `db:"created_at"`
}

func toCoreTripView(dbTripView dbTripView) trip.TripView {

	trip := trip.TripView{
		ID:                   dbTripView.ID,
		DriverName:           dbTripView.DriverName,
		DriverBrand:          dbTripView.DriverBrand,
		DriverModel:          dbTripView.DriverModel,
		DriverColor:          dbTripView.DriverColor,
		DriverPlate:          dbTripView.DriverPlate,
		SourceName:           dbTripView.SourceName,
		SourcePlaceID:        dbTripView.SourcePlaceID,
		SourceLatitude:       dbTripView.SourceLatitude,
		SourceLongitude:      dbTripView.SourceLongitude,
		DestinationName:      dbTripView.DestinationName,
		DestinationPlaceID:   dbTripView.DestinationPlaceID,
		DestinationLatitude:  dbTripView.DestinationLatitude,
		DestinationLongitude: dbTripView.DestinationLongitude,
		Status:               dbTripView.Status,
		StartTime:            dbTripView.StartTime.In(time.Local),
		CreatedAt:            dbTripView.CreatedAt.In(time.Local),
	}

	return trip
}

func toCoreTripViewSlice(dbTripViews []dbTripView) []trip.TripView {
	tripViews := make([]trip.TripView, len(dbTripViews))
	for i, dbTripView := range dbTripViews {
		tripViews[i] = toCoreTripView(dbTripView)
	}
	return tripViews
}
