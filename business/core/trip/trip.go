package trip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("trip not found")
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
	Update(ctx context.Context, trip Trip) error
	// Delete(ctx context.Context, trip Trip) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]TripView, error)
	QueryByID(ctx context.Context, tripID uuid.UUID) (TripView, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryMyTrip(ctx context.Context, userID uuid.UUID, filter QueryFilterByUser, orderBy order.By, pageNumber int, rowsPerPage int) ([]UserTrip, error)
	Join(ctx context.Context, tripPassenger TripPassenger) error

	QueryPassengers(ctx context.Context, tripID uuid.UUID) (TripDetails, error)
	UpdatePassengerStatus(ctx context.Context, tripPassenger TripPassenger) error
	CreateRating(ctx context.Context, rating Rating) error
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

type LocationRequest struct {
	Email     string `json:"email"`
	DelayTime int    `json:"delayTime"`
}

func (c *Core) PubEventToMQ(ctx context.Context, startTime time.Time, email string) error {
	// calculate the delay time
	delayTime := int(time.Until(startTime).Milliseconds())

	if delayTime < 0 {
		delayTime = 0
	}

	// send an http request post to svc tuber-location-api 3003
	// to create a new location
	locationRequest := LocationRequest{
		Email:     email,
		DelayTime: delayTime,
	}

	jsonData, err := json.Marshal(locationRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal location request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://tuber-location-api:3003/v1/tasks", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Handle non-OK response
		return fmt.Errorf("received non-OK response from location service: %s", resp.Status)
	}

	return nil
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
		UpdatedAt: now,
	}

	if err := c.storer.Create(ctx, trip); err != nil {
		return Trip{}, fmt.Errorf("create: %w", err)
	}

	return trip, nil
}

// Update replaces a user document in the database.
func (c *Core) Update(ctx context.Context, trip TripView, ut UpdateTrip) (Trip, error) {
	if ut.PassengerLimit != nil {
		trip.PassengerLimit = *ut.PassengerLimit
	}

	if ut.Status != nil {
		trip.Status = *ut.Status
	}

	trip.UpdatedAt = time.Now()

	buildTrip := Trip{
		ID:             trip.ID,
		DriverID:       trip.DriverID,
		PassengerLimit: trip.PassengerLimit,
		Status:         trip.Status,
		Source: TripLocation{
			ID: trip.SourceID,
		},
		Destination: TripLocation{
			ID: trip.DestinationID,
		},
		StartTime: trip.StartTime,
		CreatedAt: trip.CreatedAt,
		UpdatedAt: trip.UpdatedAt,
	}

	if err := c.storer.Update(ctx, buildTrip); err != nil {
		return Trip{}, fmt.Errorf("update: %w", err)
	}

	return buildTrip, nil
}

// Delete removes a user from the database.
// func (c *Core) Delete(ctx context.Context, trip User) error {
// 	if err := c.storer.Delete(ctx, trip); err != nil {
// 		return fmt.Errorf("delete: %w", err)
// 	}

// 	return nil
// }

// QueryQuery retrieves a list of existing trips from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]TripView, error) {
	trips, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return trips, nil
}

// Count returns the total number of users in the store.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// QueryMyTrip returns the trip with the specified userID from the database.
func (c *Core) QueryMyTrip(ctx context.Context, userID uuid.UUID, filter QueryFilterByUser, orderBy order.By, pageNumber int, rowsPerPage int) ([]UserTrip, error) {
	trips, err := c.storer.QueryMyTrip(ctx, userID, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return trips, nil
}

// QueryByID gets the specified user from the database.
func (c *Core) QueryByID(ctx context.Context, tripID uuid.UUID) (TripView, error) {
	trip, err := c.storer.QueryByID(ctx, tripID)
	if err != nil {
		return TripView{}, fmt.Errorf("query: tripID[%s]: %w", tripID, err)
	}

	return trip, nil
}

// CreateTripPassenger inserts a new trip into the database.
func (c *Core) Join(ctx context.Context, tripID uuid.UUID, ntp NewTripPassenger) (TripPassenger, error) {
	now := time.Now()

	tripPassenger := TripPassenger{
		TripID:        tripID,
		PassengerID:   ntp.PassengerID,
		SourceID:      ntp.SourceID,
		DestinationID: ntp.DestinationID,
		Status:        StatusPending,
		CreatedAt:     now,
	}

	if err := c.storer.Join(ctx, tripPassenger); err != nil {
		return TripPassenger{}, fmt.Errorf("create: %w", err)
	}

	return tripPassenger, nil
}

func (c *Core) UpdatePassengerStatus(ctx context.Context, tripID uuid.UUID, passengerID uuid.UUID, status string) (TripPassenger, error) {
	tripPassenger := TripPassenger{
		TripID:      tripID,
		PassengerID: passengerID,
		Status:      status,
	}

	if err := c.storer.UpdatePassengerStatus(ctx, tripPassenger); err != nil {
		return TripPassenger{}, fmt.Errorf("create: %w", err)
	}

	return tripPassenger, nil
}

// QueryPassengers retrieves a list of existing trips from the database.
func (c *Core) QueryPassengers(ctx context.Context, tripID uuid.UUID) (TripDetails, error) {
	tripDetails, err := c.storer.QueryPassengers(ctx, tripID)
	if err != nil {
		return TripDetails{}, fmt.Errorf("query: %w", err)
	}

	return tripDetails, nil
}

// CreateRating inserts a new rating into the database.
func (c *Core) CreateRating(ctx context.Context, tripID uuid.UUID, nr NewRating) (Rating, error) {
	now := time.Now()

	rating := Rating{
		ID:          uuid.New(),
		TripID:      tripID,
		CommenterID: nr.CommenterID,
		Comment:     nr.Comment,
		Rating:      nr.Rating,
		CreatedAt:   now,
	}

	if err := c.storer.CreateRating(ctx, rating); err != nil {
		return Rating{}, fmt.Errorf("create: %w", err)
	}

	return rating, nil
}
