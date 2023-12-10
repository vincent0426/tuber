package location

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	ctx := context.TODO()
	newLocation := NewLocation{
		Name:    "Test Location",
		PlaceID: "Place123",
		Lat:     10.0,
		Lon:     20.0,
	}

	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	expectedLocation := Location{
		ID:      uuid.New(),
		Name:    newLocation.Name,
		PlaceID: newLocation.PlaceID,
		Lat:     newLocation.Lat,
		Lon:     newLocation.Lon,
	}

	mockStorer.On("Create", mock.Anything, mock.AnythingOfType("Location")).Return(nil)

	location, err := core.Create(ctx, newLocation)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocation.Name, location.Name)
	mockStorer.AssertExpectations(t)
}

func TestCreateError(t *testing.T) {
	ctx := context.TODO()
	newLocation := NewLocation{
		Name:    "Test Location",
		PlaceID: "Place123",
		Lat:     10.0,
		Lon:     20.0,
	}

	mockStorer := new(MockStorer)
	mockStorer.On("Create", ctx, mock.AnythingOfType("Location")).Return(errors.New("create error"))

	core := NewCore(mockStorer)

	_, err := core.Create(ctx, newLocation)
	assert.Error(t, err)
	fmt.Println(err.Error())
	assert.Contains(t, err.Error(), "create error")
	mockStorer.AssertExpectations(t)
}

func TestQuery(t *testing.T) {
	ctx := context.TODO()
	filter := QueryFilter{} // Define filter according to your needs
	orderBy := order.By{}   // Define order according to your needs
	pageNumber := 1
	rowsPerPage := 10

	expectedLocations := []Location{
		{ID: uuid.New(), Name: "Location 1", PlaceID: "PlaceID1", Lat: 10.0, Lon: 20.0},
		{ID: uuid.New(), Name: "Location 2", PlaceID: "PlaceID2", Lat: 15.0, Lon: 25.0},
	}

	mockStorer := new(MockStorer)
	mockStorer.On("Query", ctx, filter, orderBy, pageNumber, rowsPerPage).Return(expectedLocations, nil)

	core := NewCore(mockStorer)

	locations, err := core.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocations, locations)
	mockStorer.AssertExpectations(t)
}

func TestQueryError(t *testing.T) {
	ctx := context.TODO()
	filter := QueryFilter{}
	orderBy := order.By{}
	pageNumber := 1
	rowsPerPage := 10

	mockStorer := new(MockStorer)
	mockStorer.On("Query", ctx, filter, orderBy, pageNumber, rowsPerPage).Return([]Location{}, errors.New("query error"))

	core := NewCore(mockStorer)

	_, err := core.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query error")
	mockStorer.AssertExpectations(t)
}

func TestCount(t *testing.T) {
	ctx := context.TODO()
	filter := QueryFilter{} // Define filter according to your needs

	expectedCount := 5

	mockStorer := new(MockStorer)
	mockStorer.On("Count", ctx, filter).Return(expectedCount, nil)

	core := NewCore(mockStorer)

	count, err := core.Count(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockStorer.AssertExpectations(t)
}

func TestCountError(t *testing.T) {
	ctx := context.TODO()
	filter := QueryFilter{}

	mockStorer := new(MockStorer)
	mockStorer.On("Count", ctx, filter).Return(0, errors.New("count error"))

	core := NewCore(mockStorer)

	_, err := core.Count(ctx, filter)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "count error")
	mockStorer.AssertExpectations(t)
}

func TestQueryByID(t *testing.T) {
	ctx := context.TODO()
	locationID := uuid.New()

	expectedLocation := Location{ID: locationID, Name: "Location", PlaceID: "PlaceID", Lat: 10.0, Lon: 20.0}

	mockStorer := new(MockStorer)
	mockStorer.On("QueryByID", ctx, locationID).Return(expectedLocation, nil)

	core := NewCore(mockStorer)

	location, err := core.QueryByID(ctx, locationID)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocation, location)
	mockStorer.AssertExpectations(t)
}

func TestQueryByIDError(t *testing.T) {
	ctx := context.TODO()
	locationID := uuid.New()

	mockStorer := new(MockStorer)
	mockStorer.On("QueryByID", ctx, locationID).Return(Location{}, errors.New("query by id error"))

	core := NewCore(mockStorer)

	_, err := core.QueryByID(ctx, locationID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query by id error")
	mockStorer.AssertExpectations(t)
}
