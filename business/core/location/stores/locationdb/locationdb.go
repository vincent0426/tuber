// Package locationdb contains user related CRUD functionality.
package locationdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/TSMC-Uber/server/business/core/location"
	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/TSMC-Uber/server/business/sys/database"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/google/uuid"
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

// Create inserts a new location into the database.
func (s *Store) Create(ctx context.Context, location location.Location) error {
	dbLocation := toDBLocation(location)

	sql, args, err := sq.
		Insert("locations").
		Columns("id", "name", "place_id", "lat_lon").
		Values(dbLocation.ID, dbLocation.Name, dbLocation.PlaceID, squirrel.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", dbLocation.Lat, dbLocation.Lon)).
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

// // Update replaces a user document in the database.
// func (s *Store) Update(ctx context.Context, usr user.User) error {
// 	dbUser := toDBUser(usr)

// 	sql, args, err := sq.
// 		Update("users").
// 		Set("name", dbUser.Name).
// 		Set("email", dbUser.Email).
// 		Set("bio", dbUser.Bio).
// 		Set("accept_notification", dbUser.AcceptNotification).
// 		Set("updated_at", dbUser.UpdatedAt).
// 		Where(sq.Eq{"id": dbUser.ID}).
// 		PlaceholderFormat(sq.Dollar).
// 		ToSql()

// 	if err != nil {
// 		return fmt.Errorf("tosql: %w", err)
// 	}

// 	// execute the sql
// 	if err := database.ExecContext(ctx, s.log, s.db, sql, args); err != nil {
// 		if errors.Is(err, database.ErrDBDuplicatedEntry) {
// 			return user.ErrUniqueEmail
// 		}
// 		return fmt.Errorf("execcontext: %w", err)
// 	}

// 	return nil
// }

// // Delete removes a user from the database.
// func (s *Store) Delete(ctx context.Context, usr user.User) error {
// 	sql, args, err := sq.
// 		Delete("users").
// 		Where(sq.Eq{"id": usr.ID}).
// 		PlaceholderFormat(sq.Dollar).
// 		ToSql()

// 	if err != nil {
// 		return fmt.Errorf("tosql: %w", err)
// 	}

// 	// execute the sql
// 	if err := database.ExecContext(ctx, s.log, s.db, sql, args); err != nil {
// 		return fmt.Errorf("execcontext: %w", err)
// 	}

// 	return nil
// }

// Query retrieves a list of existing locations from the database.
func (s *Store) Query(ctx context.Context, filter location.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]location.Location, error) {
	builder := sq.Select("id", "name", "place_id", "ST_Y(lat_lon::geometry) AS lat", "ST_X(lat_lon::geometry) AS lon").From("locations")

	builder = s.applyFilter(builder, filter)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
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

	var dbLocations []dbLocation
	if err := database.QueryContext(ctx, s.log, s.db, sql, nil, &dbLocations); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreLocationSlice(dbLocations), nil
}

// Count returns the total number of trips in the DB.
func (s *Store) Count(ctx context.Context, filter location.QueryFilter) (int, error) {
	builder := sq.Select("COUNT(*) AS count").From("locations")

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
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified trip from the database.
func (s *Store) QueryByID(ctx context.Context, locationID uuid.UUID) (location.Location, error) {
	sql, args, err := sq.
		Select("*").
		From("locations").
		Where(sq.Eq{"id": locationID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return location.Location{}, fmt.Errorf("tosql: %w", err)
	}

	var dbLocation dbLocation
	if err := database.GetContext(ctx, s.log, s.db, sql, args, &dbLocation); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return location.Location{}, fmt.Errorf("getcontext: %w", location.ErrNotFound)
		}
		return location.Location{}, fmt.Errorf("getcontext: %w", err)
	}

	return toCoreLocation(dbLocation), nil
}
