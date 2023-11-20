// Package userdb contains user related CRUD functionality.
package tripdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/TSMC-Uber/server/business/core/trip"
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

func insertLocationAndGetID(ctx context.Context, tx *sql.Tx, loc dbLocation) (uuid.UUID, error) {
	var locationID uuid.UUID

	// Prepare statement to insert a new location and return its ID
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO locations (name, place_id, lat_lon) VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326)) RETURNING id")
	if err != nil {
		return locationID, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	// Execute the prepared statement
	err = stmt.QueryRowContext(ctx, loc.Name, loc.PlaceID, loc.Lon, loc.Lat).Scan(&locationID)
	if err != nil {
		return locationID, fmt.Errorf("queryrow: %w", err)
	}

	return locationID, nil
}

// Create inserts a new trip into the database.
func (s *Store) Create(ctx context.Context, trip trip.Trip) error {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	// Defer to handle transaction commit/rollback
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	// Insert source and destination locations and get their IDs
	sourceID, err := insertLocationAndGetID(ctx, tx, toDBLocation(trip.Source))
	if err != nil {
		return fmt.Errorf("insert source: %w", err)
	}

	destinationID, err := insertLocationAndGetID(ctx, tx, toDBLocation(trip.Destination))
	if err != nil {
		return fmt.Errorf("insert destination: %w", err)
	}

	// Insert trip with source and destination IDs
	sql, args, err := sq.
		Insert("trip").
		Columns("id", "driver_id", "passenger_limit", "source_id", "destination_id", "status", "start_time", "created_at").
		Values(trip.ID, trip.DriverID, trip.PassengerLimit, sourceID, destinationID, trip.Status, trip.StartTime, trip.CreatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("tosql: %w", err)
	}
	if _, err := tx.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	// Insert mid locations and associate them with the trip
	for _, mid := range trip.Mid {
		midID, err := insertLocationAndGetID(ctx, tx, toDBLocation(mid))
		if err != nil {
			return fmt.Errorf("insert mid: %w", err)
		}

		// Insert into trip_location table
		sql, args, err := sq.
			Insert("trip_location").
			Columns("trip_id", "location_id").
			Values(trip.ID, midID).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("tosql: %w", err)
		}
		if _, err := tx.ExecContext(ctx, sql, args...); err != nil {
			return fmt.Errorf("execcontext: %w", err)
		}
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

// // QueryAll retrieves a list of existing trips from the database.
func (s *Store) QueryAll(ctx context.Context, filter trip.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]trip.TripView, error) {
	builder := sq.Select(
		"trip_view.id",
		"trip_view.driver_name",
		"trip_view.driver_image_url",
		"trip_view.driver_brand",
		"trip_view.driver_model",
		"trip_view.driver_color",
		"trip_view.driver_plate",
		"trip_view.source_name",
		"trip_view.source_place_id",
		"ST_Y(trip_view.source_lat_lon::geometry) AS source_latitude",
		"ST_X(trip_view.source_lat_lon::geometry) AS source_longitude",
		"trip_view.destination_name",
		"trip_view.destination_place_id",
		"ST_Y(trip_view.destination_lat_lon::geometry) AS destination_latitude",
		"ST_X(trip_view.destination_lat_lon::geometry) AS destination_longitude",
		"trip_view.status",
		"trip_view.start_time",
		"trip_view.created_at",
	).From("trip_view")

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

	var dbTripViews []dbTripView
	if err := database.QueryContext(ctx, s.log, s.db, sql, nil, &dbTripViews); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreTripViewSlice(dbTripViews), nil
}

// Query retrieves a trip from the database.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID, filter trip.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]trip.UserTrip, error) {
	sql, args, err := sq.
		Select(
			"trip_passenger.passenger_id",
			"trip_passenger.source_id",
			"trip_passenger.destination_id",
			"trip_passenger.status AS tp_status",
			"trip.id AS trip_id",
			"trip.driver_id",
			"trip.passenger_limit",
			"trip.source_id",
			"trip.destination_id",
			"trip.status AS trip_status",
			"trip.start_time",
			"trip.created_at AS trip_created_at",
		).
		From("trip_passenger").
		Join("trip ON trip_passenger.trip_id = trip.id").
		Where(sq.Eq{"trip_passenger.passenger_id": userID.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("tosql: %w", err)
	}

	var dbTrips []dbUserTrip
	if err := database.QueryContext(ctx, s.log, s.db, sql, args, &dbTrips); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreUserTripSlice(dbTrips), nil
}

func (s *Store) CreateTripPassenger(ctx context.Context, tripPassenger trip.TripPassenger) error {
	dbTripPassenger := toDBTripPassenger(tripPassenger)
	sql, args, err := sq.
		Insert("trip_passenger").
		Columns("trip_id", "passenger_id", "source_id", "destination_id", "status", "created_at").
		Values(dbTripPassenger.TripID, dbTripPassenger.PassengerID, dbTripPassenger.SourceID, dbTripPassenger.DestinationID, dbTripPassenger.Status, dbTripPassenger.CreatedAt).
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

// // Count returns the total number of trips in the DB.
func (s *Store) Count(ctx context.Context, filter trip.QueryFilter) (int, error) {
	builder := sq.Select("COUNT(*) AS count").From("trip")

	builder = s.applyFilter(builder, filter)

	// Convert the builder to SQL and args
	sql, _, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("tosql: %w", err)
	}

	var count struct {
		Count int `db:"count"`
	}
	if err := database.GetContext(ctx, s.log, s.db, sql, nil, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified trip from the database.
func (s *Store) QueryByID(ctx context.Context, tripID string) (trip.TripView, error) {
	sql, args, err := sq.Select(
		"trip_view.id",
		"trip_view.driver_name",
		"trip_view.driver_image_url",
		"trip_view.driver_brand",
		"trip_view.driver_model",
		"trip_view.driver_color",
		"trip_view.driver_plate",
		"trip_view.source_name",
		"trip_view.source_place_id",
		"ST_Y(trip_view.source_lat_lon::geometry) AS source_latitude",
		"ST_X(trip_view.source_lat_lon::geometry) AS source_longitude",
		"trip_view.destination_name",
		"trip_view.destination_place_id",
		"ST_Y(trip_view.destination_lat_lon::geometry) AS destination_latitude",
		"ST_X(trip_view.destination_lat_lon::geometry) AS destination_longitude",
		"trip_view.status",
		"trip_view.start_time",
		"trip_view.created_at",
	).From("trip_view").
		Where(sq.Eq{"id": tripID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return trip.TripView{}, fmt.Errorf("tosql: %w", err)
	}

	var dbTrip dbTripView
	if err := database.GetContext(ctx, s.log, s.db, sql, args, &dbTrip); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return trip.TripView{}, fmt.Errorf("namedquerystruct: %w", trip.ErrNotFound)
		}
		return trip.TripView{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreTripView(dbTrip), nil
}

// // QueryByIDs gets the specified users from the database.
// func (s *Store) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]user.User, error) {
// 	ids := make([]string, len(userIDs))
// 	for i, userID := range userIDs {
// 		ids[i] = userID.String()
// 	}

// 	data := struct {
// 		UserID interface {
// 			driver.Valuer
// 			sql.Scanner
// 		} `db:"user_id"`
// 	}{
// 		UserID: dbarray.Array(ids),
// 	}

// 	const q = `
// 	SELECT
// 		*
// 	FROM
// 		users
// 	WHERE
// 		user_id = ANY(:user_id)`

// 	var usrs []dbUser
// 	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &usrs); err != nil {
// 		if errors.Is(err, database.ErrDBNotFound) {
// 			return nil, user.ErrNotFound
// 		}
// 		return nil, fmt.Errorf("namedquerystruct: %w", err)
// 	}

// 	return toCoreUserSlice(usrs), nil
// }

// // QueryByEmail gets the specified user from the database by email.
// func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
// 	data := struct {
// 		Email string `db:"email"`
// 	}{
// 		Email: email.Address,
// 	}

// 	const q = `
// 	SELECT
// 		*
// 	FROM
// 		users
// 	WHERE
// 		email = :email`

// 	var dbUsr dbUser
// 	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
// 		if errors.Is(err, database.ErrDBNotFound) {
// 			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
// 		}
// 		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
// 	}

// 	return toCoreUser(dbUsr), nil
// }

// // QueryByGoogleID gets the specified user from the database by googleID.
// func (s *Store) QueryByGoogleID(ctx context.Context, googleID string) (user.User, error) {
// 	sql, args, err := sq.
// 		Select("*").
// 		From("users").
// 		Where(sq.Eq{"sub": googleID}).
// 		PlaceholderFormat(sq.Dollar).
// 		ToSql()

// 	if err != nil {
// 		return user.User{}, fmt.Errorf("tosql: %w", err)
// 	}

// 	var dbUsr dbUser
// 	if err := database.GetContext(ctx, s.log, s.db, sql, args, &dbUsr); err != nil {
// 		if errors.Is(err, database.ErrDBNotFound) {
// 			return user.User{}, fmt.Errorf("namedquerystruct: %w", database.ErrDBNotFound)
// 		}
// 		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
// 	}

// 	return toCoreUser(dbUsr), nil
// }
