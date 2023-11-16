package locationdb

import (
	"github.com/TSMC-Uber/server/business/core/location"
	"github.com/google/uuid"
)

type dbLocation struct {
	ID      uuid.UUID `db:"id"`
	Name    string    `db:"name"`
	PlaceID string    `db:"place_id"`
	Lat     float64   `db:"lat"`
	Lon     float64   `db:"lon"`
}

func toDBLocation(location location.Location) dbLocation {
	return dbLocation{
		ID:      location.ID,
		Name:    location.Name,
		PlaceID: location.PlaceID,
		Lat:     location.Lat,
		Lon:     location.Lon,
	}
}

func toCoreLocation(dbLocation dbLocation) location.Location {

	trip := location.Location{
		ID:      dbLocation.ID,
		Name:    dbLocation.Name,
		PlaceID: dbLocation.PlaceID,
		Lat:     dbLocation.Lat,
		Lon:     dbLocation.Lon,
	}

	return trip
}

func toCoreLocationSlice(dbLocations []dbLocation) []location.Location {
	locations := make([]location.Location, len(dbLocations))
	for i, dbLocation := range dbLocations {
		locations[i] = toCoreLocation(dbLocation)
	}
	return locations
}
