// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/TSMC-Uber/server/business/sys/database"
	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
	QueryByGoogleID(ctx context.Context, googleID string) (User, error)
}

// Core manages the set of APIs for user access.
type Core struct {
	storer Storer
}

// NewCore constructs a core for user api access.
func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

// Create inserts a new user into the database.
func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	now := time.Now()

	usr := User{
		ID:                 uuid.New(),
		Name:               nu.Name,
		Email:              nu.Email,
		Bio:                nu.Bio,
		AcceptNotification: nu.AcceptNotification,
		Sub:                nu.Sub,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := c.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// Update replaces a user document in the database.
func (c *Core) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Bio != nil {
		usr.Bio = *uu.Bio
	}
	if uu.AcceptNotification != nil {
		usr.AcceptNotification = *uu.AcceptNotification
	}

	usr.UpdatedAt = time.Now()

	if err := c.storer.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	return usr, nil
}

// Delete removes a user from the database.
func (c *Core) Delete(ctx context.Context, usr User) error {
	if err := c.storer.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error) {
	users, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

// Count returns the total number of users in the store.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// QueryByID gets the specified user from the database.
func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := c.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
	}

	return user, nil
}

// QueryByIDs gets the specified user from the database.
func (c *Core) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error) {
	user, err := c.storer.QueryByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("query: userIDs[%s]: %w", userIDs, err)
	}

	return user, nil
}

// QueryByEmail gets the specified user from the database by email.
func (c *Core) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	user, err := c.storer.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	return user, nil
}

// UpsertByGoogleID gets the specified user from the database by email.
func (c *Core) UpsertByGoogleID(ctx context.Context, googleID string, nu NewUser) (User, error) {
	usr, err := c.storer.QueryByGoogleID(ctx, googleID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrDBNotFound):
			usr, err = c.Create(ctx, nu)
			if err != nil {
				return User{}, fmt.Errorf("create: %w", err)
			}
		default:
			return User{}, fmt.Errorf("query: googleID[%s]: %w", googleID, err)
		}
	}

	return usr, nil
}
