package driver

import (
	"context"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockStorer struct {
	mock.Mock
}

func (m *MockStorer) Create(ctx context.Context, driver Driver) error {
	args := m.Called(ctx, driver)
	return args.Error(0)
}

func (m *MockStorer) QueryByID(ctx context.Context, driverID string) (Driver, error) {
	args := m.Called(ctx, driverID)
	return args.Get(0).(Driver), args.Error(1)
}

func (m *MockStorer) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Driver, error) {
	args := m.Called(ctx, filter, orderBy, pageNumber, rowsPerPage)
	return args.Get(0).([]Driver), args.Error(1)
}

func (m *MockStorer) Count(ctx context.Context, filter QueryFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func (m *MockStorer) AddFavorite(ctx context.Context, userID uuid.UUID, driverID string) error {
	args := m.Called(ctx, userID, driverID)
	return args.Error(0)
}

func (m *MockStorer) QueryFavorite(ctx context.Context, userID uuid.UUID, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]FavoriteDriver, error) {
	args := m.Called(ctx, userID, filter, orderBy, pageNumber, rowsPerPage)
	return args.Get(0).([]FavoriteDriver), args.Error(1)
}
