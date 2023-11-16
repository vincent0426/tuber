package location

import (
	"context"
	"errors"
	"fmt"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound = errors.New("location not found")

// ErrUniqueEmail           = errors.New("email is not unique")
// ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, location Location) error
	// Update(ctx context.Context, trip Trip) error
	// Delete(ctx context.Context, trip Trip) error
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Location, error)
	QueryByID(ctx context.Context, locationID string) (Location, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
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
func (c *Core) Create(ctx context.Context, nu NewLocation) (Location, error) {

	location := Location{
		ID:      uuid.New(),
		Name:    nu.Name,
		PlaceID: nu.PlaceID,
		Lat:     nu.Lat,
		Lon:     nu.Lon,
	}
	fmt.Println("core: trip: create: location:", location)
	if err := c.storer.Create(ctx, location); err != nil {
		return Location{}, fmt.Errorf("create: %w", err)
	}

	return location, nil
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
func (c *Core) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Location, error) {
	locations, err := c.storer.QueryAll(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return locations, nil
}

// // Count returns the total number of users in the store.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// // QueryByID gets the specified user from the database.
func (c *Core) QueryByID(ctx context.Context, locationID string) (Location, error) {
	location, err := c.storer.QueryByID(ctx, locationID)
	if err != nil {
		return Location{}, fmt.Errorf("query: tripID[%s]: %w", locationID, err)
	}

	return location, nil
}

// QueryByIDs gets the specified user from the database.
// func (c *Core) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]Trip, error) {
// 	user, err := c.storer.QueryByIDs(ctx, userIDs)
// 	if err != nil {
// 		return nil, fmt.Errorf("query: userIDs[%s]: %w", userIDs, err)
// 	}

// 	return user, nil
// }
