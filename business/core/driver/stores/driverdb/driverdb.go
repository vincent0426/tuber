// Package driverdb contains driver related CRUD functionality.
package driverdb

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/TSMC-Uber/server/business/core/driver"
	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/TSMC-Uber/server/business/sys/database"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  *sqlx.DB
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new trip into the database.
func (s *Store) Create(ctx context.Context, driver driver.Driver) error {
	dbDriver := toDBDriver(driver)
	fmt.Println("store: driver: create: dbDriver:", dbDriver)
	sql, args, err := sq.
		Insert("driver").
		Columns("user_id", "license", "verified", "brand", "model", "color", "plate", "created_at").
		Values(dbDriver.UserID, dbDriver.License, dbDriver.Verified, dbDriver.Brand, dbDriver.Model, dbDriver.Color, dbDriver.Plate, dbDriver.CreatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("tosql: %w", err)
	}

	// execute the sql
	if err := database.ExecContext(ctx, s.log, s.db, sql, args); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

// QueryAll retrieves a list of existing drivers from the database.
func (s *Store) QueryAll(ctx context.Context, filter driver.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]driver.Driver, error) {
	builder := sq.Select(
		"id", // user_id
		"name",
		"image_url",
		"license",
		"verified",
		"brand",
		"model",
		"color",
		"plate",
		"driver_created_at",
	).From("driver_view")

	builder = s.applyFilter(builder, filter)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		fmt.Println("store: trip: queryall: err:", err)
		return nil, err
	}
	builder = builder.OrderBy(orderByClause)

	// add paging
	builder = builder.Limit(uint64(rowsPerPage)).Offset(uint64((pageNumber - 1) * rowsPerPage))

	// Convert the builder to SQL and args
	sql, _, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("tosql: %w", err)
	}

	var dbDrivers []dbDriver
	if err := database.QueryContext(ctx, s.log, s.db, sql, nil, &dbDrivers); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreDriverSlice(dbDrivers), nil
}

// Count returns the total number of drivers in the DB.
func (s *Store) Count(ctx context.Context, filter driver.QueryFilter) (int, error) {
	builder := sq.Select("COUNT(*) AS count").From("driver_view")

	builder = s.applyFilter(builder, filter)

	// Convert the builder to SQL and args
	sql, _, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("tosql: %w", err)
	}

	var count struct {
		Count int `db:"count"`
	}
	if err = database.GetContext(ctx, s.log, s.db, sql, nil, &count); err != nil {
		return 0, fmt.Errorf("getcontext: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified trip from the database.
func (s *Store) QueryByID(ctx context.Context, driverID string) (driver.Driver, error) {
	sql, args, err := sq.Select(
		"id", // user_id
		"name",
		"image_url",
		"license",
		"verified",
		"brand",
		"model",
		"color",
		"plate",
		"driver_created_at",
	).From("driver_view").
		Where(sq.Eq{"id": driverID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return driver.Driver{}, fmt.Errorf("tosql: %w", err)
	}

	var dbDriver dbDriver
	if err := database.GetContext(ctx, s.log, s.db, sql, args, &dbDriver); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return driver.Driver{}, fmt.Errorf("getcontext: %w", driver.ErrNotFound)
		}
		return driver.Driver{}, fmt.Errorf("getcontext: %w", err)
	}

	return toCoreDriver(dbDriver), nil
}
