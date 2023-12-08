package user

import (
	"context"
	"net/mail"
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

	nu := NewUser{
		Name:               "Test User",
		Email:              mail.Address{Address: "000.gmail.com"},
		ImageURL:           "https://example.com/image.png",
		Bio:                "This is a test user.",
		AcceptNotification: true,
		Sub:                "1234567890",
	}

	mockStorer.On("Create", context.Background(), mock.AnythingOfType("User")).Return(nil)
	usr, err := core.Create(context.Background(), nu)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, usr.ID)
	assert.Equal(t, nu.Name, usr.Name)
	assert.Equal(t, nu.Email, usr.Email)
	assert.Equal(t, nu.ImageURL, usr.ImageURL)
	assert.Equal(t, nu.Bio, usr.Bio)
	assert.Equal(t, nu.AcceptNotification, usr.AcceptNotification)
	assert.Equal(t, nu.Sub, usr.Sub)
	assert.NotEqual(t, uuid.Nil, usr.CreatedAt)
	assert.NotEqual(t, uuid.Nil, usr.UpdatedAt)

	mockStorer.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	userID := uuid.New()
	user := User{
		ID:                 userID,
		Name:               "Test User",
		Email:              mail.Address{Address: "000.gmail.com"},
		ImageURL:           "https://example.com/image.png",
		Bio:                "This is a test user.",
		AcceptNotification: true,
		Sub:                "1234567890",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	updatedName := "Updated Name"
	updatedEmail := mail.Address{Address: "111.gmail.com"}
	updatedImageURL := "https://example.com/image2.png"
	updatedBio := "This is an updated test user."
	updatedAcceptNotification := false
	updateUser := UpdateUser{
		Name:               &updatedName,
		Email:              &updatedEmail,
		ImageURL:           &updatedImageURL,
		Bio:                &updatedBio,
		AcceptNotification: &updatedAcceptNotification,
	}

	mockStorer.On("Update", mock.Anything, mock.AnythingOfType("User")).Return(nil)
	updatedUser, err := core.Update(context.Background(), user, updateUser)

	assert.NoError(t, err)
	assert.Equal(t, updatedName, updatedUser.Name)
	assert.Equal(t, updatedEmail, updatedUser.Email)
	assert.Equal(t, updatedImageURL, updatedUser.ImageURL)
	assert.Equal(t, updatedBio, updatedUser.Bio)
	assert.Equal(t, updatedAcceptNotification, updatedUser.AcceptNotification)
	mockStorer.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	usr := User{
		ID:                 uuid.New(),
		Name:               "Test User",
		Email:              mail.Address{Address: "000.gmail.com"},
		ImageURL:           "https://example.com/image.png",
		Bio:                "This is a test user.",
		AcceptNotification: true,
		Sub:                "1234567890",
	}

	mockStorer.On("Delete", context.Background(), usr).Return(nil)
	err := core.Delete(context.Background(), usr)
	assert.NoError(t, err)

	mockStorer.AssertExpectations(t)
}

func TestQuery(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	filter := QueryFilter{
		Name:  strPtr("Test User"),
		Email: &mail.Address{Address: "000.gmail.com"},
	}

	orderBy := order.By{}

	pageNumber := 1
	rowsPerPage := 10

	mockStorer.On("Query", context.Background(), filter, orderBy, pageNumber, rowsPerPage).Return([]User{}, nil)
	_, err := core.Query(context.Background(), filter, orderBy, pageNumber, rowsPerPage)
	assert.NoError(t, err)

	mockStorer.AssertExpectations(t)
}

func TestCount(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	filter := QueryFilter{
		Name:  strPtr("Test User"),
		Email: &mail.Address{Address: "000.gmail.com"},
	}

	mockStorer.On("Count", context.Background(), filter).Return(0, nil)
	_, err := core.Count(context.Background(), filter)
	assert.NoError(t, err)

	mockStorer.AssertExpectations(t)
}

func TestQueryByID(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	userID := uuid.New()

	mockStorer.On("QueryByID", context.Background(), userID).Return(User{}, nil)
	_, err := core.QueryByID(context.Background(), userID)
	assert.NoError(t, err)

	mockStorer.AssertExpectations(t)
}

func TestQueryByIDs(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	userIDs := []uuid.UUID{uuid.New(), uuid.New()}

	mockStorer.On("QueryByIDs", context.Background(), userIDs).Return([]User{}, nil)
	_, err := core.QueryByIDs(context.Background(), userIDs)
	assert.NoError(t, err)

	mockStorer.AssertExpectations(t)
}

func TestQueryByEmail(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	email := mail.Address{Address: "000.gmail.com"}

	mockStorer.On("QueryByEmail", context.Background(), email).Return(User{}, nil)
	_, err := core.QueryByEmail(context.Background(), email)
	assert.NoError(t, err)

	mockStorer.AssertExpectations(t)
}

func strPtr(s string) *string {
	return &s
}

func TestUpsertByGoogleID(t *testing.T) {
	mockStorer := new(MockStorer)
	core := NewCore(mockStorer)

	newUser := NewUser{
		Name:               "Test User",
		Email:              mail.Address{Address: "000.gmail.com"},
		ImageURL:           "https://example.com/image.png",
		Bio:                "This is a test user.",
		AcceptNotification: true,
		Sub:                "1234567890",
	}
	googleID := "1234567890"

	mockStorer.On("QueryByGoogleID", context.Background(), googleID).Return(User{}, nil)
	_, err := core.UpsertByGoogleID(context.Background(), googleID, newUser)
	assert.NoError(t, err)
	mockStorer.AssertExpectations(t)
}
