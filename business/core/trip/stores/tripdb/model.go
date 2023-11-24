package tripdb

import (
	"time"

	"github.com/TSMC-Uber/server/business/core/trip"
	"github.com/google/uuid"
)

type dbTrip struct {
	ID             uuid.UUID `db:"id"`
	DriverID       uuid.UUID `db:"driver_id"`
	PassengerLimit int       `db:"passenger_limit"`
	SourceID       uuid.UUID `db:"source_id"`
	DestinationID  uuid.UUID `db:"destination_id"`
	Status         string    `db:"status"`
	StartTime      time.Time `db:"start_time"`
	CreatedAt      time.Time `db:"trip_created_at"`
	UpdatedAt      time.Time `db:"trip_updated_at"`
}

func toDBTrip(trip trip.Trip) dbTrip {
	return dbTrip{
		ID:             trip.ID,
		DriverID:       trip.DriverID,
		PassengerLimit: trip.PassengerLimit,
		SourceID:       trip.Source.ID,
		DestinationID:  trip.Destination.ID,
		Status:         trip.Status,
		StartTime:      trip.StartTime.UTC(),
		CreatedAt:      trip.CreatedAt.UTC(),
		UpdatedAt:      trip.UpdatedAt.UTC(),
	}
}

// ------------------------------------------------------------
type dbUserTrip struct {
	TripID          uuid.UUID `db:"trip_id"`
	PassengerID     uuid.UUID `db:"passenger_id"`
	MySourceID      uuid.UUID `db:"my_source_id"`
	MyDestinationID uuid.UUID `db:"my_destination_id"`
	MyStatus        string    `db:"my_status"`
	DriverID        uuid.UUID `db:"driver_id"`
	PassengerLimit  int       `db:"passenger_limit"`
	SourceID        uuid.UUID `db:"source_id"`
	DestinationID   uuid.UUID `db:"destination_id"`
	TripStatus      string    `db:"trip_status"`
	StartTime       time.Time `db:"start_time"`
	CreatedAt       time.Time `db:"trip_created_at"`
	UpdatedAt       time.Time `db:"trip_updated_at"`
}

func toCoreUserTrip(dbUserTrip dbUserTrip) trip.UserTrip {

	trip := trip.UserTrip{
		TripID:          dbUserTrip.TripID,
		PassengerID:     dbUserTrip.PassengerID,
		MySourceID:      dbUserTrip.MySourceID,
		MyDestinationID: dbUserTrip.MyDestinationID,
		MyStatus:        dbUserTrip.MyStatus,
		DriverID:        dbUserTrip.DriverID,
		PassengerLimit:  dbUserTrip.PassengerLimit,
		SourceID:        dbUserTrip.SourceID,
		DestinationID:   dbUserTrip.DestinationID,
		TripStatus:      dbUserTrip.TripStatus,
		StartTime:       dbUserTrip.StartTime.In(time.Local),
		CreatedAt:       dbUserTrip.CreatedAt.In(time.Local),
		UpdatedAt:       dbUserTrip.UpdatedAt.In(time.Local),
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
	DriverID             uuid.UUID `db:"driver_id"`
	DriverName           string    `db:"driver_name"`
	DriverImageURL       string    `db:"driver_image_url"`
	DriverBrand          string    `db:"driver_brand"`
	DriverModel          string    `db:"driver_model"`
	DriverColor          string    `db:"driver_color"`
	DriverPlate          string    `db:"driver_plate"`
	SourceID             uuid.UUID `db:"source_id"`
	SourceName           string    `db:"source_name"`
	SourcePlaceID        string    `db:"source_place_id"`
	SourceLatitude       float64   `db:"source_latitude"`
	SourceLongitude      float64   `db:"source_longitude"`
	DestinationID        uuid.UUID `db:"destination_id"`
	DestinationName      string    `db:"destination_name"`
	DestinationPlaceID   string    `db:"destination_place_id"`
	DestinationLatitude  float64   `db:"destination_latitude"`
	DestinationLongitude float64   `db:"destination_longitude"`
	PassengerLimit       int       `db:"passenger_limit"`
	Status               string    `db:"status"`
	StartTime            time.Time `db:"start_time"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}

func toCoreTripView(dbTripView dbTripView) trip.TripView {

	trip := trip.TripView{
		ID:                   dbTripView.ID,
		DriverID:             dbTripView.DriverID,
		DriverName:           dbTripView.DriverName,
		DriverImageURL:       dbTripView.DriverImageURL,
		DriverBrand:          dbTripView.DriverBrand,
		DriverModel:          dbTripView.DriverModel,
		DriverColor:          dbTripView.DriverColor,
		DriverPlate:          dbTripView.DriverPlate,
		SourceID:             dbTripView.SourceID,
		SourceName:           dbTripView.SourceName,
		SourcePlaceID:        dbTripView.SourcePlaceID,
		SourceLatitude:       dbTripView.SourceLatitude,
		SourceLongitude:      dbTripView.SourceLongitude,
		DestinationID:        dbTripView.DestinationID,
		DestinationName:      dbTripView.DestinationName,
		DestinationPlaceID:   dbTripView.DestinationPlaceID,
		DestinationLatitude:  dbTripView.DestinationLatitude,
		DestinationLongitude: dbTripView.DestinationLongitude,
		PassengerLimit:       dbTripView.PassengerLimit,
		Status:               dbTripView.Status,
		StartTime:            dbTripView.StartTime.In(time.Local),
		CreatedAt:            dbTripView.CreatedAt.In(time.Local),
		UpdatedAt:            dbTripView.UpdatedAt.In(time.Local),
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

// ------------------------------------------------------------
type dbLocation struct {
	Name    string  `db:"name"`
	PlaceID string  `db:"place_id"`
	Lat     float64 `db:"lat"`
	Lon     float64 `db:"lon"`
}

func toDBLocation(location trip.TripLocation) dbLocation {
	return dbLocation{
		Name:    location.Name,
		PlaceID: location.PlaceID,
		Lat:     location.Lat,
		Lon:     location.Lon,
	}
}

func toCoreTripDetails(tripDetails trip.TripDetails) trip.TripDetails {
	trip := trip.TripDetails{
		TripID:               tripDetails.TripID,
		DriverID:             tripDetails.DriverID,
		DriverName:           tripDetails.DriverName,
		DriverImageURL:       tripDetails.DriverImageURL,
		DriverBrand:          tripDetails.DriverBrand,
		DriverModel:          tripDetails.DriverModel,
		DriverColor:          tripDetails.DriverColor,
		DriverPlate:          tripDetails.DriverPlate,
		SourceName:           tripDetails.SourceName,
		SourcePlaceID:        tripDetails.SourcePlaceID,
		SourceLatitude:       tripDetails.SourceLatitude,
		SourceLongitude:      tripDetails.SourceLongitude,
		DestinationName:      tripDetails.DestinationName,
		DestinationPlaceID:   tripDetails.DestinationPlaceID,
		DestinationLatitude:  tripDetails.DestinationLatitude,
		DestinationLongitude: tripDetails.DestinationLongitude,
		PassengerDetails:     tripDetails.PassengerDetails,
	}

	return trip
}

type dbRating struct {
	ID          uuid.UUID `db:"id"`
	TripID      uuid.UUID `db:"trip_id"`
	CommenterID uuid.UUID `db:"commenter_id"`
	Comment     string    `db:"comment"`
	Rating      int       `db:"rating"`
	CreatedAt   time.Time `db:"created_at"`
}

func toDBRating(rating trip.Rating) dbRating {
	return dbRating{
		ID:          rating.ID,
		TripID:      rating.TripID,
		CommenterID: rating.CommenterID,
		Comment:     rating.Comment,
		Rating:      rating.Rating,
		CreatedAt:   rating.CreatedAt.UTC(),
	}
}
