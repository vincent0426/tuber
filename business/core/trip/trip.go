package trip

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TSMC-Uber/server/business/data/order"
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
	Create(ctx context.Context, trip Trip) error
	// Update(ctx context.Context, trip Trip) error
	// Delete(ctx context.Context, trip Trip) error
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Trip, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]UserTrip, error)
	QueryByID(ctx context.Context, tripID uuid.UUID) (Trip, error)
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

// Create inserts a new trip into the database.
func (c *Core) Create(ctx context.Context, nu NewTrip) (Trip, error) {
	now := time.Now()

	trip := Trip{
		ID:             uuid.New(),
		DriverID:       nu.DriverID,
		PassengerLimit: nu.PassengerLimit,
		SourceID:       nu.SourceID,
		DestinationID:  nu.DestinationID,
		StartTime:      nu.StartTime,
		CreatedAt:      now,
	}
	fmt.Println("core: trip: create: trip:", trip)
	if err := c.storer.Create(ctx, trip); err != nil {
		return Trip{}, fmt.Errorf("create: %w", err)
	}

	return trip, nil
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
func (c *Core) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Trip, error) {
	trips, err := c.storer.QueryAll(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return trips, nil
}

// Count returns the total number of users in the store.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// QueryByUserID returns the trip with the specified userID from the database.
func (c *Core) QueryByUserID(ctx context.Context, userID uuid.UUID, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]UserTrip, error) {
	fmt.Println("core: trip: querybyuserid: userID:", userID)
	trips, err := c.storer.QueryByUserID(ctx, userID, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	fmt.Println("core: trip: querybyuserid: trips:", trips)
	return trips, nil
}

// QueryByID gets the specified user from the database.
func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (Trip, error) {
	user, err := c.storer.QueryByID(ctx, userID)
	if err != nil {
		return Trip{}, fmt.Errorf("query: tripID[%s]: %w", userID, err)
	}

	return user, nil
}

// QueryByIDs gets the specified user from the database.
// func (c *Core) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]Trip, error) {
// 	user, err := c.storer.QueryByIDs(ctx, userIDs)
// 	if err != nil {
// 		return nil, fmt.Errorf("query: userIDs[%s]: %w", userIDs, err)
// 	}

// 	return user, nil
// }
