package trip

import (
	"context"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockStorer struct {
	mock.Mock
}

func (m *MockStorer) Create(ctx context.Context, trip Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockStorer) Update(ctx context.Context, trip Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockStorer) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]TripView, error) {
	args := m.Called(ctx, filter, orderBy, pageNumber, rowsPerPage)
	return args.Get(0).([]TripView), args.Error(1)
}

func (m *MockStorer) QueryByID(ctx context.Context, tripID uuid.UUID) (TripView, error) {
	args := m.Called(ctx, tripID)
	return args.Get(0).(TripView), args.Error(1)
}

func (m *MockStorer) Count(ctx context.Context, filter QueryFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func (m *MockStorer) QueryMyTrip(ctx context.Context, userID uuid.UUID, filter QueryFilterByUser, orderBy order.By, pageNumber int, rowsPerPage int) ([]UserTrip, error) {
	args := m.Called(ctx, userID, filter, orderBy, pageNumber, rowsPerPage)
	return args.Get(0).([]UserTrip), args.Error(1)
}

func (m *MockStorer) Join(ctx context.Context, tripPassenger TripPassenger) error {
	args := m.Called(ctx, tripPassenger)
	return args.Error(0)
}

func (m *MockStorer) QueryPassengers(ctx context.Context, tripID uuid.UUID) (TripDetails, error) {
	args := m.Called(ctx, tripID)
	return args.Get(0).(TripDetails), args.Error(1)
}

func (m *MockStorer) CreateRating(ctx context.Context, rating Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}
