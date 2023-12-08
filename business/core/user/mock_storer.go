package user

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockStorer struct {
	mock.Mock
}

func (m *MockStorer) Create(ctx context.Context, usr User) error {
	fmt.Println("MockStorer.Create", usr)
	args := m.Called(ctx, usr)
	return args.Error(0)
}

func (m *MockStorer) Update(ctx context.Context, usr User) error {
	fmt.Println("MockStorer.Update", usr)
	args := m.Called(ctx, usr)
	return args.Error(0)
}

func (m *MockStorer) Delete(ctx context.Context, usr User) error {
	fmt.Println("MockStorer.Delete", usr)
	args := m.Called(ctx, usr)
	return args.Error(0)
}

func (m *MockStorer) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error) {
	fmt.Println("MockStorer.Query", filter)
	args := m.Called(ctx, filter, orderBy, pageNumber, rowsPerPage)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockStorer) Count(ctx context.Context, filter QueryFilter) (int, error) {
	fmt.Println("MockStorer.Count", filter)
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func (m *MockStorer) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	fmt.Println("MockStorer.QueryByID", userID)
	args := m.Called(ctx, userID)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockStorer) QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error) {
	fmt.Println("MockStorer.QueryByIDs", userID)
	args := m.Called(ctx, userID)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockStorer) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	fmt.Println("MockStorer.QueryByEmail", email)
	args := m.Called(ctx, email)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockStorer) QueryByGoogleID(ctx context.Context, googleID string) (User, error) {
	fmt.Println("MockStorer.QueryByGoogleID", googleID)
	args := m.Called(ctx, googleID)
	return args.Get(0).(User), args.Error(1)
}
