// Package userdb contains user related CRUD functionality.
package tripdb

import (
	"context"
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

// // Create inserts a new trip into the database.
func (s *Store) Create(ctx context.Context, trip trip.Trip) error {
	dbTrip := toDBTrip(trip)
	fmt.Println("store: trip: create: dbTrip:", dbTrip)
	sql, args, err := sq.
		Insert("trip").
		Columns("id", "driver_id", "passenger_limit", "source_id", "destination_id", "start_time", "created_at").
		Values(dbTrip.ID, dbTrip.DriverID, dbTrip.PassengerLimit, dbTrip.SourceID, dbTrip.DestinationID, dbTrip.StartTime, dbTrip.CreatedAt).
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

// // QueryAll retrieves a list of existing trips from the database.
func (s *Store) QueryAll(ctx context.Context, filter trip.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]trip.Trip, error) {
	builder := sq.Select("*").From("trip")

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

	var dbTrips []dbTrip
	if err := database.SelectContext(ctx, s.log, s.db, sql, nil, &dbTrips); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreTripSlice(dbTrips), nil
}

// Query retrieves a trip from the database.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID, filter trip.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]trip.UserTrip, error) {
	sql, args, err := sq.
		Select(
			"trip_passenger.passenger_id",
			"trip_passenger.station_source_id",
			"trip_passenger.station_destination_id",
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
	if err := database.Select(ctx, s.log, s.db, sql, args, &dbTrips); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreUserTripSlice(dbTrips), nil
}

// // Count returns the total number of trips in the DB.
func (s *Store) Count(ctx context.Context, filter trip.QueryFilter) (int, error) {
	// data := map[string]interface{}{}

	// builder := sq.Select("COUNT(*) AS count").From("users")

	// builder = s.applyFilter(builder, filter)

	// // Convert the builder to SQL and args
	// sql, _, err := builder.ToSql()
	// if err != nil {
	// 	return 0, fmt.Errorf("tosql: %w", err)
	// }

	// var count struct {
	// 	Count int `db:"count"`
	// }
	// if err := database.NamedQueryStruct(ctx, s.log, s.db, sql, data, &count); err != nil {
	// 	return 0, fmt.Errorf("namedquerystruct: %w", err)
	// }

	return 0, nil
}

// // QueryByID gets the specified trip from the database.
func (s *Store) QueryByID(ctx context.Context, tripID uuid.UUID) (trip.Trip, error) {
	// sql, args, err := sq.
	// 	Select("*").
	// 	From("users").
	// 	Where(sq.Eq{"id": userID}).
	// 	PlaceholderFormat(sq.Dollar).
	// 	ToSql()

	// if err != nil {
	// 	return user.User{}, fmt.Errorf("tosql: %w", err)
	// }

	// var dbUsr dbUser
	// if err := database.GetContext(ctx, s.log, s.db, sql, args, &dbUsr); err != nil {
	// 	if errors.Is(err, database.ErrDBNotFound) {
	// 		return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
	// 	}
	// 	return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	// }

	// return toCoreUser(dbUsr), nil
	return trip.Trip{}, nil
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
