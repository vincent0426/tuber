package driver

import (
	"context"
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

	newDriver := NewDriver{
		UserID:  uuid.New(),
		License: "some-license",
		Brand:   "some-brand",
		Model:   "some-model",
		Color:   "some-color",
		Plate:   "some-plate",
	}

	expectedDriver := Driver{
		UserID:    newDriver.UserID,
		License:   newDriver.License,
		Verified:  DriverNotVerified,
		Brand:     newDriver.Brand,
		Model:     newDriver.Model,
		Color:     newDriver.Color,
		Plate:     newDriver.Plate,
		CreatedAt: time.Now(),
	}
	mockStorer.On("Create", mock.Anything, mock.AnythingOfType("Driver")).Return(nil)

	createdDriver, err := core.Create(context.Background(), newDriver)

	assert.NoError(t, err)
	assert.Equal(t, expectedDriver.UserID, createdDriver.UserID)
	assert.Equal(t, expectedDriver.License, createdDriver.License)
	assert.Equal(t, expectedDriver.Verified, createdDriver.Verified)
	assert.Equal(t, expectedDriver.Brand, createdDriver.Brand)
	assert.Equal(t, expectedDriver.Model, createdDriver.Model)
	assert.Equal(t, expectedDriver.Color, createdDriver.Color)
	assert.Equal(t, expectedDriver.Plate, createdDriver.Plate)
	mockStorer.AssertExpectations(t)
}

func TestQuery(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	filter := QueryFilter{} // Set appropriate filter
	orderBy := order.By{}   // Set appropriate order
	pageNumber := 1
	rowsPerPage := 10

	expectedDrivers := []Driver{
		{UserID: uuid.New(), License: "123", Verified: DriverVerified, Brand: "1-brand", Model: "1-model", Color: "1-color", Plate: "1-plate"},
		{UserID: uuid.New(), License: "456", Verified: DriverVerified, Brand: "2-brand", Model: "2-model", Color: "2-color", Plate: "2-plate"},
		{UserID: uuid.New(), License: "789", Verified: DriverVerified, Brand: "3-brand", Model: "3-model", Color: "3-color", Plate: "3-plate"},
	}
	ctx := context.Background()
	mockStorer.On("Query", ctx, filter, orderBy, pageNumber, rowsPerPage).Return(expectedDrivers, nil)

	drivers, err := core.Query(context.Background(), filter, orderBy, pageNumber, rowsPerPage)

	assert.NoError(t, err)
	assert.Equal(t, expectedDrivers, drivers)
	mockStorer.AssertExpectations(t)
}

func TestQueryByID(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	driverID := "some-id"
	expectedDriver := Driver{
		UserID:    uuid.New(),
		License:   "some-license",
		Verified:  DriverVerified,
		Brand:     "some-brand",
		Model:     "some-model",
		Color:     "some-color",
		Plate:     "some-plate",
		CreatedAt: time.Now(),
	}

	mockStorer.On("QueryByID", mock.Anything, driverID).Return(expectedDriver, nil)

	driver, err := core.QueryByID(context.Background(), driverID)

	assert.NoError(t, err)
	assert.Equal(t, expectedDriver, driver)
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

func TestAddFavorite(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	userID := uuid.New()
	driverID := "some-id"

	mockStorer.On("AddFavorite", mock.Anything, userID, driverID).Return(nil)

	err := core.AddFavorite(context.Background(), userID, driverID)

	assert.NoError(t, err)
	mockStorer.AssertExpectations(t)
}

func TestQueryFavorite(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	userID := uuid.New()
	filter := QueryFilter{} // Set appropriate filter
	orderBy := order.By{}   // Set appropriate order
	pageNumber := 1
	rowsPerPage := 10

	expectedFavoriteDrivers := []FavoriteDriver{
		{DriverID: uuid.New(), DriverName: "1-name", DriverImageURL: "1-image-url", DriverBrand: "1-brand", DriverModel: "1-model", DriverColor: "1-color", DriverPlate: "1-plate"},
		{DriverID: uuid.New(), DriverName: "2-name", DriverImageURL: "2-image-url", DriverBrand: "2-brand", DriverModel: "2-model", DriverColor: "2-color", DriverPlate: "2-plate"},
		{DriverID: uuid.New(), DriverName: "3-name", DriverImageURL: "3-image-url", DriverBrand: "3-brand", DriverModel: "3-model", DriverColor: "3-color", DriverPlate: "3-plate"},
	}
	ctx := context.Background()
	mockStorer.On("QueryFavorite", ctx, userID, filter, orderBy, pageNumber, rowsPerPage).Return(expectedFavoriteDrivers, nil)

	favoriteDrivers, err := core.QueryFavorite(context.Background(), userID, filter, orderBy, pageNumber, rowsPerPage)

	assert.NoError(t, err)
	assert.Equal(t, expectedFavoriteDrivers, favoriteDrivers)
	mockStorer.AssertExpectations(t)
}
