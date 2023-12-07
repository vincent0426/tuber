package trip

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	newTrip := NewTrip{
		DriverID:       uuid.New(),
		PassengerLimit: 4,
		Source: TripLocation{
			Name:    "Test Location",
			PlaceID: "Place123",
			Lat:     10.0,
			Lon:     20.0,
		},
		Destination: TripLocation{
			Name:    "Test Location",
			PlaceID: "Place123",
			Lat:     10.0,
			Lon:     20.0,
		},
		Mid: []TripLocation{
			{
				Name:    "Test Location",
				PlaceID: "Place123",
				Lat:     10.0,
				Lon:     20.0,
			},
		},
		StartTime: time.Now(),
	}

	expectedTrip := TripView{
		ID:                   uuid.New(),
		DriverID:             newTrip.DriverID,
		PassengerLimit:       newTrip.PassengerLimit,
		SourceID:             uuid.New(),
		SourceName:           newTrip.Source.Name,
		SourcePlaceID:        newTrip.Source.PlaceID,
		SourceLatitude:       newTrip.Source.Lat,
		SourceLongitude:      newTrip.Source.Lon,
		DestinationID:        uuid.New(),
		DestinationName:      newTrip.Destination.Name,
		DestinationPlaceID:   newTrip.Destination.PlaceID,
		DestinationLatitude:  newTrip.Destination.Lat,
		DestinationLongitude: newTrip.Destination.Lon,
		Mid:                  newTrip.Mid,
		Status:               TripStatusNotStarted,
		StartTime:            newTrip.StartTime,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	mockStorer.On("Create", mock.Anything, mock.AnythingOfType("Trip")).Return(nil)

	createdTrip, err := core.Create(context.Background(), newTrip)

	assert.NoError(t, err)
	assert.Equal(t, expectedTrip.DriverID, createdTrip.DriverID)
	assert.Equal(t, expectedTrip.PassengerLimit, createdTrip.PassengerLimit)
	assert.Equal(t, expectedTrip.SourceName, createdTrip.Source.Name)
	assert.Equal(t, expectedTrip.SourcePlaceID, createdTrip.Source.PlaceID)
	assert.Equal(t, expectedTrip.SourceLatitude, createdTrip.Source.Lat)
	assert.Equal(t, expectedTrip.SourceLongitude, createdTrip.Source.Lon)
	assert.Equal(t, expectedTrip.DestinationName, createdTrip.Destination.Name)
	assert.Equal(t, expectedTrip.DestinationPlaceID, createdTrip.Destination.PlaceID)
	assert.Equal(t, expectedTrip.DestinationLatitude, createdTrip.Destination.Lat)
	assert.Equal(t, expectedTrip.DestinationLongitude, createdTrip.Destination.Lon)
	assert.Equal(t, expectedTrip.Mid, createdTrip.Mid)
	assert.Equal(t, expectedTrip.Status, createdTrip.Status)
	assert.Equal(t, expectedTrip.StartTime, createdTrip.StartTime)
	mockStorer.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	tripID := uuid.New()
	trip := TripView{
		ID:                   tripID,
		DriverID:             uuid.New(),
		PassengerLimit:       4,
		SourceID:             uuid.New(),
		SourceName:           "Test Location",
		SourcePlaceID:        "Place123",
		SourceLatitude:       10.0,
		SourceLongitude:      20.0,
		DestinationID:        uuid.New(),
		DestinationName:      "Test Location",
		DestinationPlaceID:   "Place123",
		DestinationLatitude:  10.0,
		DestinationLongitude: 20.0,
		Mid: []TripLocation{
			{
				ID:      uuid.New(),
				Name:    "Test Location",
				PlaceID: "Place123",
				Lat:     10.0,
				Lon:     20.0,
			},
		},
		Status:    TripStatusNotStarted,
		StartTime: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	updatedPassengerLimit := 4
	updatedStatus := TripStatusNotStarted

	updateTrip := UpdateTrip{
		PassengerLimit: &updatedPassengerLimit,
		Status:         &updatedStatus,
	}

	expectedTrip := TripView{
		PassengerLimit: updatedPassengerLimit,
		Status:         updatedStatus,
	}
	mockStorer.On("Update", mock.Anything, mock.AnythingOfType("Trip")).Return(nil)
	updatedTrip, err := core.Update(context.Background(), trip, updateTrip)

	assert.NoError(t, err)
	assert.Equal(t, expectedTrip.PassengerLimit, updatedTrip.PassengerLimit)
	assert.Equal(t, expectedTrip.Status, updatedTrip.Status)
	mockStorer.AssertExpectations(t)
}

func TestQuery(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	filter := QueryFilter{} // Set appropriate filter
	orderBy := order.By{}   // Set appropriate order
	pageNumber := 1
	rowsPerPage := 10

	expectedTrips := []TripView{
		{
			ID:                   uuid.New(),
			DriverID:             uuid.New(),
			PassengerLimit:       4,
			SourceID:             uuid.New(),
			SourceName:           "Test Location",
			SourcePlaceID:        "Place123",
			SourceLatitude:       10.0,
			SourceLongitude:      20.0,
			DestinationID:        uuid.New(),
			DestinationName:      "Test Location",
			DestinationPlaceID:   "Place123",
			DestinationLatitude:  10.0,
			DestinationLongitude: 20.0,
			Mid: []TripLocation{
				{
					ID:      uuid.New(),
					Name:    "Test Location",
					PlaceID: "Place123",
					Lat:     10.0,
					Lon:     20.0,
				},
			},
			Status:    TripStatusNotStarted,
			StartTime: time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:                   uuid.New(),
			DriverID:             uuid.New(),
			PassengerLimit:       4,
			SourceID:             uuid.New(),
			SourceName:           "Test Location",
			SourcePlaceID:        "Place123",
			SourceLatitude:       10.0,
			SourceLongitude:      20.0,
			DestinationID:        uuid.New(),
			DestinationName:      "Test Location",
			DestinationPlaceID:   "Place123",
			DestinationLatitude:  10.0,
			DestinationLongitude: 20.0,
			Mid: []TripLocation{
				{
					ID:      uuid.New(),
					Name:    "Test Location",
					PlaceID: "Place123",
					Lat:     10.0,
					Lon:     20.0,
				},
			},
			Status:    TripStatusNotStarted,
			StartTime: time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockStorer.On("Query", mock.Anything, filter, orderBy, pageNumber, rowsPerPage).Return(expectedTrips, nil)

	trips, err := core.Query(context.Background(), filter, orderBy, pageNumber, rowsPerPage)

	assert.NoError(t, err)
	for i, trip := range trips {
		assert.Equal(t, expectedTrips[i].DriverID, trip.DriverID)
		assert.Equal(t, expectedTrips[i].PassengerLimit, trip.PassengerLimit)
		assert.Equal(t, expectedTrips[i].SourceName, trip.SourceName)
		assert.Equal(t, expectedTrips[i].SourcePlaceID, trip.SourcePlaceID)
		assert.Equal(t, expectedTrips[i].SourceLatitude, trip.SourceLatitude)
		assert.Equal(t, expectedTrips[i].SourceLongitude, trip.SourceLongitude)
		assert.Equal(t, expectedTrips[i].DestinationName, trip.DestinationName)
		assert.Equal(t, expectedTrips[i].DestinationPlaceID, trip.DestinationPlaceID)
		assert.Equal(t, expectedTrips[i].DestinationLatitude, trip.DestinationLatitude)
		assert.Equal(t, expectedTrips[i].DestinationLongitude, trip.DestinationLongitude)
		assert.Equal(t, expectedTrips[i].Mid, trip.Mid)
		assert.Equal(t, expectedTrips[i].Status, trip.Status)
		assert.Equal(t, expectedTrips[i].StartTime, trip.StartTime)
	}
	mockStorer.AssertExpectations(t)
}

func TestQueryByID(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	tripID := uuid.New()

	expectedTrip := TripView{
		ID:                   tripID,
		DriverID:             uuid.New(),
		PassengerLimit:       4,
		SourceID:             uuid.New(),
		SourceName:           "Test Location",
		SourcePlaceID:        "Place123",
		SourceLatitude:       10.0,
		SourceLongitude:      20.0,
		DestinationID:        uuid.New(),
		DestinationName:      "Test Location",
		DestinationPlaceID:   "Place123",
		DestinationLatitude:  10.0,
		DestinationLongitude: 20.0,
		Mid: []TripLocation{
			{
				ID:      uuid.New(),
				Name:    "Test Location",
				PlaceID: "Place123",
				Lat:     10.0,
				Lon:     20.0,
			},
		},
		Status:    TripStatusNotStarted,
		StartTime: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockStorer.On("QueryByID", mock.Anything, tripID).Return(expectedTrip, nil)

	trip, err := core.QueryByID(context.Background(), tripID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTrip.DriverID, trip.DriverID)
	assert.Equal(t, expectedTrip.PassengerLimit, trip.PassengerLimit)
	assert.Equal(t, expectedTrip.SourceName, trip.SourceName)
	assert.Equal(t, expectedTrip.SourcePlaceID, trip.SourcePlaceID)
	assert.Equal(t, expectedTrip.SourceLatitude, trip.SourceLatitude)
	assert.Equal(t, expectedTrip.SourceLongitude, trip.SourceLongitude)
	assert.Equal(t, expectedTrip.DestinationName, trip.DestinationName)
	assert.Equal(t, expectedTrip.DestinationPlaceID, trip.DestinationPlaceID)
	assert.Equal(t, expectedTrip.DestinationLatitude, trip.DestinationLatitude)
	assert.Equal(t, expectedTrip.DestinationLongitude, trip.DestinationLongitude)
	assert.Equal(t, expectedTrip.Mid, trip.Mid)
	assert.Equal(t, expectedTrip.Status, trip.Status)
	assert.Equal(t, expectedTrip.StartTime, trip.StartTime)
	mockStorer.AssertExpectations(t)
}

func TestQueryByIDError(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	tripID := uuid.New()

	mockStorer.On("QueryByID", mock.Anything, tripID).Return(TripView{}, errors.New("query by id error"))

	_, err := core.QueryByID(context.Background(), tripID)
	assert.Error(t, err)
	fmt.Println(err.Error())
	assert.Contains(t, err.Error(), "query by id error")
	mockStorer.AssertExpectations(t)
}

func TestCount(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	filter := QueryFilter{} // Set appropriate filter

	expectedCount := 10
	mockStorer.On("Count", mock.Anything, filter).Return(expectedCount, nil)

	count, err := core.Count(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockStorer.AssertExpectations(t)
}

func TestCountError(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	filter := QueryFilter{} // Set appropriate filter

	mockStorer.On("Count", mock.Anything, filter).Return(0, errors.New("count error"))

	_, err := core.Count(context.Background(), filter)
	assert.Error(t, err)
	fmt.Println(err.Error())
	assert.Contains(t, err.Error(), "count error")
	mockStorer.AssertExpectations(t)
}

func TestQueryMyTrip(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	userID := uuid.New()
	filter := QueryFilterByUser{} // Set appropriate filter
	orderBy := order.By{}         // Set appropriate order
	pageNumber := 1
	rowsPerPage := 10

	expectedTrips := []UserTrip{
		{
			TripID:          uuid.New(),
			PassengerID:     uuid.New(),
			MySourceID:      uuid.New(),
			MyDestinationID: uuid.New(),
			MyStatus:        TripStatusNotStarted,
			DriverID:        uuid.New(),
			PassengerLimit:  4,
			SourceID:        uuid.New(),
			DestinationID:   uuid.New(),
			TripStatus:      TripStatusNotStarted,
			StartTime:       time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			TripID:          uuid.New(),
			PassengerID:     uuid.New(),
			MySourceID:      uuid.New(),
			MyDestinationID: uuid.New(),
			MyStatus:        TripStatusNotStarted,
			DriverID:        uuid.New(),
			PassengerLimit:  4,
			SourceID:        uuid.New(),
			DestinationID:   uuid.New(),
			TripStatus:      TripStatusNotStarted,
			StartTime:       time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	mockStorer.On("QueryMyTrip", mock.Anything, userID, filter, orderBy, pageNumber, rowsPerPage).Return(expectedTrips, nil)

	trips, err := core.QueryMyTrip(context.Background(), userID, filter, orderBy, pageNumber, rowsPerPage)

	assert.NoError(t, err)
	for i, trip := range trips {
		assert.Equal(t, expectedTrips[i].TripID, trip.TripID)
		assert.Equal(t, expectedTrips[i].PassengerID, trip.PassengerID)
		assert.Equal(t, expectedTrips[i].MySourceID, trip.MySourceID)
		assert.Equal(t, expectedTrips[i].MyDestinationID, trip.MyDestinationID)
		assert.Equal(t, expectedTrips[i].MyStatus, trip.MyStatus)
		assert.Equal(t, expectedTrips[i].DriverID, trip.DriverID)
		assert.Equal(t, expectedTrips[i].PassengerLimit, trip.PassengerLimit)
		assert.Equal(t, expectedTrips[i].SourceID, trip.SourceID)
		assert.Equal(t, expectedTrips[i].DestinationID, trip.DestinationID)
		assert.Equal(t, expectedTrips[i].TripStatus, trip.TripStatus)
		assert.Equal(t, expectedTrips[i].StartTime, trip.StartTime)
	}
	mockStorer.AssertExpectations(t)
}
