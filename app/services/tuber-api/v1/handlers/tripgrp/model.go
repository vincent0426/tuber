package tripgrp

import (
	"time"

	"github.com/TSMC-Uber/server/business/core/trip"
	"github.com/TSMC-Uber/server/business/sys/validate"
	"github.com/google/uuid"
)

// AppTrip represents information about an individual trip.
type AppTrip struct {
	ID             string `json:"id"`
	DriverID       string `json:"driver_id"`
	PassengerLimit int    `json:"passenger_limit"`
	SourceID       string `json:"source_id"`
	DestinationID  string `json:"destination_id"`
	Status         string `json:"status"`
	StartTime      string `json:"start_time"`
	CreatedAt      string `json:"createdAt"`
}

func toAppTrip(trip trip.Trip) AppTrip {

	return AppTrip{
		ID:             trip.ID.String(),
		DriverID:       trip.DriverID.String(),
		PassengerLimit: trip.PassengerLimit,
		SourceID:       trip.Source.ID.String(),
		DestinationID:  trip.Destination.ID.String(),
		Status:         trip.Status,
		StartTime:      trip.StartTime.Format(time.RFC3339),
		CreatedAt:      trip.CreatedAt.Format(time.RFC3339),
	}
}

type AppNewTripLocation struct {
	Name    string  `json:"name" binding:"required"`
	PlaceID string  `json:"place_id" binding:"required"`
	Lat     float64 `json:"lat" binding:"required"`
	Lon     float64 `json:"lon" binding:"required"`
}

type AppNewTrip struct {
	PassengerLimit int                  `json:"passenger_limit"`
	Source         AppNewTripLocation   `json:"source"`
	Destination    AppNewTripLocation   `json:"destination"`
	Mid            []AppNewTripLocation `json:"mid"`
	StartTime      string               `json:"start_time" binding:"required"`
}

func toCoreNewTrip(app AppNewTrip) (trip.NewTrip, error) {
	// turn string to time
	startTime, err := time.Parse(time.RFC3339, app.StartTime)
	if err != nil {
		return trip.NewTrip{}, err
	}

	// construct mid
	mid := []trip.TripLocation{}
	for _, appMid := range app.Mid {
		mid = append(mid, trip.TripLocation{
			Name:    appMid.Name,
			PlaceID: appMid.PlaceID,
			Lat:     appMid.Lat,
			Lon:     appMid.Lon,
		})
	}

	trip := trip.NewTrip{
		PassengerLimit: app.PassengerLimit,
		Source: trip.TripLocation{
			Name:    app.Source.Name,
			PlaceID: app.Source.PlaceID,
			Lat:     app.Source.Lat,
			Lon:     app.Source.Lon,
		},
		Destination: trip.TripLocation{
			Name:    app.Destination.Name,
			PlaceID: app.Destination.PlaceID,
			Lat:     app.Destination.Lat,
			Lon:     app.Destination.Lon,
		},
		Mid:       mid,
		StartTime: startTime,
	}

	return trip, nil
}

type AppUpdateTrip struct {
	PassengerLimit int    `json:"passenger_limit"`
	Status         string `json:"status"`
}

func toCoreUpdateTrip(app AppUpdateTrip) (trip.UpdateTrip, error) {

	trip := trip.UpdateTrip{
		PassengerLimit: &app.PassengerLimit,
		Status:         &app.Status,
	}

	return trip, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewTrip) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

type AppTripPassenger struct {
	TripID        string `json:"trip_id"`
	PassengerID   string `json:"passenger_id"`
	SourceID      string `json:"source_id"`
	DestinationID string `json:"destination_id"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
}

func toAppTripPassenger(tripPassenger trip.TripPassenger) AppTripPassenger {

	return AppTripPassenger{
		TripID:        tripPassenger.TripID.String(),
		PassengerID:   tripPassenger.PassengerID.String(),
		SourceID:      tripPassenger.SourceID.String(),
		DestinationID: tripPassenger.DestinationID.String(),
		Status:        tripPassenger.Status,
		CreatedAt:     tripPassenger.CreatedAt.Format(time.RFC3339),
	}
}

// =============================================================================
type AppNewTripPassenger struct {
	SourceID      string `json:"source_id" binding:"required"`
	DestinationID string `json:"destination_id" binding:"required"`
}

func toCoreNewTripPassenger(app AppNewTripPassenger) (trip.NewTripPassenger, error) {
	uuSourceID, err := uuid.Parse(app.SourceID)
	if err != nil {
		return trip.NewTripPassenger{}, err
	}
	uuDestinationID, err := uuid.Parse(app.DestinationID)
	if err != nil {
		return trip.NewTripPassenger{}, err
	}

	tripPassenger := trip.NewTripPassenger{
		SourceID:      uuSourceID,
		DestinationID: uuDestinationID,
	}

	return tripPassenger, nil
}

type AppTripView struct {
	ID                   string            `json:"id"`
	DriverID             string            `json:"driver_id"`
	DriverName           string            `json:"driver_name"`
	DriverImageURL       string            `json:"driver_image_url"`
	DriverBrand          string            `json:"driver_brand"`
	DriverModel          string            `json:"driver_model"`
	DriverColor          string            `json:"driver_color"`
	DriverPlate          string            `json:"driver_plate"`
	SourceID             string            `json:"source_id"`
	SourceName           string            `json:"source_name"`
	SourcePlaceID        string            `json:"source_place_id"`
	SourceLatitude       float64           `json:"source_latitude"`
	SourceLongitude      float64           `json:"source_longitude"`
	DestinationID        string            `json:"destination_id"`
	DestinationName      string            `json:"destination_name"`
	DestinationPlaceID   string            `json:"destination_place_id"`
	DestinationLatitude  float64           `json:"destination_latitude"`
	DestinationLongitude float64           `json:"destination_longitude"`
	PassengerLimit       int               `json:"passenger_limit"`
	Status               string            `json:"status"`
	StartTime            string            `json:"start_time"`
	CreatedAt            string            `json:"createdAt"`
	UpdatedAt            string            `json:"updatedAt"`
	Mid                  []AppTripLocation `json:"mid"`
}

type AppTripLocation struct {
	Name    string
	PlaceID string
	Lat     float64
	Lon     float64
}

func toAppTripView(tripView trip.TripView) AppTripView {
	// convert mid to AppTripLocation
	mid := []AppTripLocation{}
	for _, tripLocation := range tripView.Mid {
		mid = append(mid, AppTripLocation{
			Name:    tripLocation.Name,
			PlaceID: tripLocation.PlaceID,
			Lat:     tripLocation.Lat,
			Lon:     tripLocation.Lon,
		})
	}

	return AppTripView{
		ID:                   tripView.ID.String(),
		DriverID:             tripView.DriverID.String(),
		DriverName:           tripView.DriverName,
		DriverImageURL:       tripView.DriverImageURL,
		DriverBrand:          tripView.DriverBrand,
		DriverModel:          tripView.DriverModel,
		DriverColor:          tripView.DriverColor,
		DriverPlate:          tripView.DriverPlate,
		SourceID:             tripView.SourceID.String(),
		SourceName:           tripView.SourceName,
		SourcePlaceID:        tripView.SourcePlaceID,
		SourceLatitude:       tripView.SourceLatitude,
		SourceLongitude:      tripView.SourceLongitude,
		DestinationID:        tripView.DestinationID.String(),
		DestinationName:      tripView.DestinationName,
		DestinationPlaceID:   tripView.DestinationPlaceID,
		DestinationLatitude:  tripView.DestinationLatitude,
		DestinationLongitude: tripView.DestinationLongitude,
		PassengerLimit:       tripView.PassengerLimit,
		Status:               tripView.Status,
		StartTime:            tripView.StartTime.Format(time.RFC3339),
		CreatedAt:            tripView.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            tripView.UpdatedAt.Format(time.RFC3339),
		Mid:                  mid,
	}
}

type AppPassengerDetails struct {
	PassengerID          string  `json:"passenger_id"`
	SourceName           string  `json:"source_name"`
	SourcePlaceID        string  `json:"source_place_id"`
	SourceLatitude       float64 `json:"source_latitude"`
	SourceLongitude      float64 `json:"source_longitude"`
	DestinationName      string  `json:"destination_name"`
	DestinationPlaceID   string  `json:"destination_place_id"`
	DestinationLatitude  float64 `json:"destination_latitude"`
	DestinationLongitude float64 `json:"destination_longitude"`
}
type AppTripDriverDetails struct {
	DriverID       string `json:"driver_id"`
	DriverName     string `json:"driver_name"`
	DriverImageURL string `json:"driver_image_url"`
	DriverBrand    string `json:"driver_brand"`
	DriverModel    string `json:"driver_model"`
	DriverColor    string `json:"driver_color"`
	DriverPlate    string `json:"driver_plate"`
}

type AppTripDetails struct {
	TripID               string                `json:"trip_id"`
	DriverDetails        AppTripDriverDetails  `json:"driver_details"`
	SourceName           string                `json:"source_name"`
	SourcePlaceID        string                `json:"source_place_id"`
	SourceLatitude       float64               `json:"source_latitude"`
	SourceLongitude      float64               `json:"source_longitude"`
	DestinationName      string                `json:"destination_name"`
	DestinationPlaceID   string                `json:"destination_place_id"`
	DestinationLatitude  float64               `json:"destination_latitude"`
	DestinationLongitude float64               `json:"destination_longitude"`
	PassengerDetails     []AppPassengerDetails `json:"passenger_details"`
}

func toAppPassengerDetails(passengerDetails []trip.PassengerDetails) []AppPassengerDetails {
	appPassengerDetails := []AppPassengerDetails{}
	for _, passengerDetail := range passengerDetails {
		appPassengerDetails = append(appPassengerDetails, AppPassengerDetails{
			PassengerID:          passengerDetail.PassengerID.String(),
			SourceName:           passengerDetail.SourceName,
			SourcePlaceID:        passengerDetail.SourcePlaceID,
			SourceLatitude:       passengerDetail.SourceLatitude,
			SourceLongitude:      passengerDetail.SourceLongitude,
			DestinationName:      passengerDetail.DestinationName,
			DestinationPlaceID:   passengerDetail.DestinationPlaceID,
			DestinationLatitude:  passengerDetail.DestinationLatitude,
			DestinationLongitude: passengerDetail.DestinationLongitude,
		})
	}
	return appPassengerDetails
}

func toAppTripDetails(tripDetails trip.TripDetails) AppTripDetails {

	return AppTripDetails{
		TripID: tripDetails.TripID.String(),
		DriverDetails: AppTripDriverDetails{
			DriverID:       tripDetails.DriverID.String(),
			DriverName:     tripDetails.DriverName,
			DriverImageURL: tripDetails.DriverImageURL,
			DriverBrand:    tripDetails.DriverBrand,
			DriverModel:    tripDetails.DriverModel,
			DriverColor:    tripDetails.DriverColor,
			DriverPlate:    tripDetails.DriverPlate,
		},
		SourceName:           tripDetails.SourceName,
		SourcePlaceID:        tripDetails.SourcePlaceID,
		SourceLatitude:       tripDetails.SourceLatitude,
		SourceLongitude:      tripDetails.SourceLongitude,
		DestinationName:      tripDetails.DestinationName,
		DestinationPlaceID:   tripDetails.DestinationPlaceID,
		DestinationLatitude:  tripDetails.DestinationLatitude,
		DestinationLongitude: tripDetails.DestinationLongitude,
		PassengerDetails:     toAppPassengerDetails(tripDetails.PassengerDetails),
	}
}

// =============================================================================
type AppRating struct {
	ID          string `json:"id"`
	TripID      string `json:"trip_id"`
	CommenterID string `json:"commenter_id"`
	Comment     string `json:"comment"`
	Rating      int    `json:"rating"`
	CreatedAt   string `json:"createdAt"`
}

func toAppRating(rating trip.Rating) AppRating {

	return AppRating{
		ID:          rating.ID.String(),
		TripID:      rating.TripID.String(),
		CommenterID: rating.CommenterID.String(),
		Comment:     rating.Comment,
		Rating:      rating.Rating,
		CreatedAt:   rating.CreatedAt.Format(time.RFC3339),
	}
}

type AppNewRating struct {
	Rating  int    `json:"rating" binding:"required"`
	Comment string `json:"comment"`
}

func toCoreNewRating(app AppNewRating) (trip.NewRating, error) {
	rating := trip.NewRating{
		Rating:  app.Rating,
		Comment: app.Comment,
	}

	return rating, nil
}
