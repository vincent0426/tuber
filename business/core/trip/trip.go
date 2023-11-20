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

var (
	TripStatusNotStarted = "not_start"
	TripStatusIn         = "in_trip"
	TripStatusFinished   = "finished"
)

var (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, trip Trip) error
	// Update(ctx context.Context, trip Trip) error
	// Delete(ctx context.Context, trip Trip) error
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]TripView, error)
	QueryByID(ctx context.Context, tripID string) (TripView, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]UserTrip, error)
	CreateTripPassenger(ctx context.Context, tripPassenger TripPassenger) error
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
func (c *Core) Create(ctx context.Context, nt NewTrip) (Trip, error) {
	now := time.Now()

	// add id to each element in []mid
	for i := range nt.Mid {
		nt.Mid[i].ID = uuid.New()
	}

	trip := Trip{
		ID:             uuid.New(),
		DriverID:       nt.DriverID,
		PassengerLimit: nt.PassengerLimit,
		Source: TripLocation{
			ID:      uuid.New(),
			Name:    nt.Source.Name,
			PlaceID: nt.Source.PlaceID,
			Lat:     nt.Source.Lat,
			Lon:     nt.Source.Lon,
		},
		Destination: TripLocation{
			ID:      uuid.New(),
			Name:    nt.Destination.Name,
			PlaceID: nt.Destination.PlaceID,
			Lat:     nt.Destination.Lat,
			Lon:     nt.Destination.Lon,
		},
		Mid:       nt.Mid,
		Status:    TripStatusNotStarted,
		StartTime: nt.StartTime,
		CreatedAt: now,
	}

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
func (c *Core) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]TripView, error) {
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
	trips, err := c.storer.QueryByUserID(ctx, userID, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return trips, nil
}

// QueryByID gets the specified user from the database.
func (c *Core) QueryByID(ctx context.Context, tripID string) (TripView, error) {
	trip, err := c.storer.QueryByID(ctx, tripID)
	if err != nil {
		return TripView{}, fmt.Errorf("query: tripID[%s]: %w", tripID, err)
	}

	return trip, nil
}

// CreateTripPassenger inserts a new trip into the database.
func (c *Core) Join(ctx context.Context, ntp NewTripPassenger) (TripPassenger, error) {
	now := time.Now()

	tripPassenger := TripPassenger{
		TripID:        ntp.TripID,
		PassengerID:   ntp.PassengerID,
		SourceID:      ntp.SourceID,
		DestinationID: ntp.DestinationID,
		Status:        StatusPending,
		CreatedAt:     now,
	}

	if err := c.storer.CreateTripPassenger(ctx, tripPassenger); err != nil {
		return TripPassenger{}, fmt.Errorf("create: %w", err)
	}

	return tripPassenger, nil
}
