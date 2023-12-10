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
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, location Location) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Location, error)
	QueryByID(ctx context.Context, locationID uuid.UUID) (Location, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}

// Core manages the set of APIs for user access.
type Core struct {
	storer Storer
}

// NewCore constructs a core for location api access.
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

	if err := c.storer.Create(ctx, location); err != nil {
		return Location{}, fmt.Errorf("create: %w", err)
	}

	return location, nil
}

// Query retrieves a list of existing locations from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Location, error) {
	locations, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return []Location{}, fmt.Errorf("query: %w", err)
	}

	return locations, nil
}

// Count returns the total number of locations in the store.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// QueryByID gets the specified location from the database.
func (c *Core) QueryByID(ctx context.Context, locationID uuid.UUID) (Location, error) {
	location, err := c.storer.QueryByID(ctx, locationID)
	if err != nil {
		return Location{}, fmt.Errorf("query: tripID[%s]: %w", locationID, err)
	}

	return location, nil
}
