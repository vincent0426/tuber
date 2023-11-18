package driver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TSMC-Uber/server/business/data/order"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound = errors.New("driver not found")
)

var (
	DriverNotVerified = false
	DriverVerified    = true
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, driver Driver) error
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Driver, error)
	QueryByID(ctx context.Context, driverID string) (Driver, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	// Update(ctx context.Context, trip Trip) error
	// Delete(ctx context.Context, trip Trip) error
	// QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]Trip, error)
	// QueryByEmail(ctx context.Context, email mail.Address) (Trip, error)
	// QueryByGoogleID(ctx context.Context, googleID string) (Trip, error)
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

// Create inserts a new location into the database.
func (c *Core) Create(ctx context.Context, nu NewDriver) (Driver, error) {

	driver := Driver{
		UserID:    nu.UserID,
		License:   nu.License,
		Verified:  DriverNotVerified,
		Brand:     nu.Brand,
		Model:     nu.Model,
		Color:     nu.Color,
		Plate:     nu.Plate,
		CreatedAt: time.Now(),
	}
	fmt.Println("core: driver: create: driver:", driver)
	if err := c.storer.Create(ctx, driver); err != nil {
		return Driver{}, fmt.Errorf("create: %w", err)
	}

	return driver, nil
}

// Update replaces a user document in the database.
// func (c *Core) Update(ctx context.Context, trip User, uu UpdateUser) (User, error) {
// 	if uu.Name != nil {
// 		trip.Name = *uu.Name
// 	}
// 	if uu.Email != nil {
// 		trip.Email = *uu.Email
// 	}
// 	if uu.Bio != nil {
// 		trip.Bio = *uu.Bio
// 	}
// 	if uu.AcceptNotification != nil {
// 		trip.AcceptNotification = *uu.AcceptNotification
// 	}

// 	trip.UpdatedAt = time.Now()

// 	if err := c.storer.Update(ctx, trip); err != nil {
// 		return User{}, fmt.Errorf("update: %w", err)
// 	}

// 	return trip, nil
// }

// Delete removes a user from the database.
// func (c *Core) Delete(ctx context.Context, trip User) error {
// 	if err := c.storer.Delete(ctx, trip); err != nil {
// 		return fmt.Errorf("delete: %w", err)
// 	}

// 	return nil
// }

// QueryQueryAll retrieves a list of existing trips from the database.
func (c *Core) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Driver, error) {
	drivers, err := c.storer.QueryAll(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return drivers, nil
}

// // // Count returns the total number of users in the store.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// QueryByID gets the specified user from the database.
func (c *Core) QueryByID(ctx context.Context, driverID string) (Driver, error) {
	driver, err := c.storer.QueryByID(ctx, driverID)
	if err != nil {
		return Driver{}, fmt.Errorf("query: driverID[%s]: %w", driverID, err)
	}

	return driver, nil
}

// QueryByIDs gets the specified user from the database.
// func (c *Core) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]Trip, error) {
// 	user, err := c.storer.QueryByIDs(ctx, userIDs)
// 	if err != nil {
// 		return nil, fmt.Errorf("query: userIDs[%s]: %w", userIDs, err)
// 	}

// 	return user, nil
// }
