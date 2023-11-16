package locationgrp

import (
	"github.com/TSMC-Uber/server/business/core/location"
)

// AppLocation represents information about an individual location.
type AppLocation struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	PlaceID string  `json:"place_id"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func toAppLocation(location location.Location) AppLocation {

	return AppLocation{
		ID:      location.ID.String(),
		Name:    location.Name,
		PlaceID: location.PlaceID,
		Lat:     location.Lat,
		Lon:     location.Lon,
	}
}

// =============================================================================

// AppNewUser contains information needed to create a new user.
type AppNewLocation struct {
	Name    string  `json:"name" binding:"required"`
	PlaceID string  `json:"place_id" binding:"required"`
	Lat     float64 `json:"lat" binding:"required"`
	Lon     float64 `json:"lon" binding:"required"`
}

func toCoreNewLocation(app AppNewLocation) (location.NewLocation, error) {

	location := location.NewLocation{
		Name:    app.Name,
		PlaceID: app.PlaceID,
		Lat:     app.Lat,
		Lon:     app.Lon,
	}

	return location, nil
}

// // Validate checks the data in the model is considered clean.
// func (app AppNewTrip) Validate() error {
// 	if err := validate.Check(app); err != nil {
// 		return err
// 	}
// 	return nil
// }
