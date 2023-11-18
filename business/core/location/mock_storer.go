package location

import (
	"context"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/stretchr/testify/mock"
)

type MockStorer struct {
	mock.Mock
}

func (m *MockStorer) Create(ctx context.Context, location Location) error {
	args := m.Called(ctx, location)
	return args.Error(0)
}

func (m *MockStorer) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Location, error) {
	args := m.Called(ctx, filter, orderBy, pageNumber, rowsPerPage)
	return args.Get(0).([]Location), args.Error(1)
}

func (m *MockStorer) QueryByID(ctx context.Context, locationID string) (Location, error) {
	args := m.Called(ctx, locationID)
	return args.Get(0).(Location), args.Error(1)
}

func (m *MockStorer) Count(ctx context.Context, filter QueryFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}
