// Package userdb contains user related CRUD functionality.
package userdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/mail"

	sq "github.com/Masterminds/squirrel"
	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/TSMC-Uber/server/business/sys/database/dbarray"

	database "github.com/TSMC-Uber/server/business/sys/database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) error {
	dbUser := toDBUser(usr)

	sql, args, err := sq.
		Insert("users").
		Columns("id", "name", "email", "bio", "language", "accept_notification").
		Values(dbUser.ID, dbUser.Name, dbUser.Email, dbUser.Bio, dbUser.Lang, dbUser.AcceptNotification).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("tosql: %w", err)
	}

	// execute the sql
	if err := database.ExecContext(ctx, s.log, s.db, sql, args); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return user.ErrUniqueEmail
		}
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, usr user.User) error {
	dbUser := toDBUser(usr)

	sql, args, err := sq.
		Update("users").
		Set("name", dbUser.Name).
		Set("email", dbUser.Email).
		Set("bio", dbUser.Bio).
		Set("language", dbUser.Lang).
		Set("accept_notification", dbUser.AcceptNotification).
		Where(sq.Eq{"id": dbUser.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("tosql: %w", err)
	}

	// execute the sql
	if err := database.ExecContext(ctx, s.log, s.db, sql, args); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return user.ErrUniqueEmail
		}
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr user.User) error {
	sql, args, err := sq.
		Delete("users").
		Where(sq.Eq{"id": usr.ID}).
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

// Query retrieves a list of existing users from the database.
func (s *Store) Query(ctx context.Context, filter user.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]user.User, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	builder := sq.Select("*").From("users")

	builder = s.applyFilter(builder, filter)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}
	builder = builder.OrderBy(orderByClause)

	// Convert the builder to SQL and args
	sql, _, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("tosql: %w", err)
	}

	var dbUsrs []dbUser
	if err := database.NamedQuerySlice(ctx, s.log, s.db, sql, data, &dbUsrs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreUserSlice(dbUsrs), nil
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter user.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	builder := sq.Select("COUNT(*) AS count").From("users")

	builder = s.applyFilter(builder, filter)

	// Convert the builder to SQL and args
	sql, _, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("tosql: %w", err)
	}

	var count struct {
		Count int `db:"count"`
	}
	if err := database.NamedQueryStruct(ctx, s.log, s.db, sql, data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	sql, args, err := sq.
		Select("*").
		From("users").
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return user.User{}, fmt.Errorf("tosql: %w", err)
	}

	var dbUsr dbUser
	if err := database.GetContext(ctx, s.log, s.db, sql, args, &dbUsr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUser(dbUsr), nil
}

// QueryByIDs gets the specified users from the database.
func (s *Store) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]user.User, error) {
	ids := make([]string, len(userIDs))
	for i, userID := range userIDs {
		ids[i] = userID.String()
	}

	data := struct {
		UserID interface {
			driver.Valuer
			sql.Scanner
		} `db:"user_id"`
	}{
		UserID: dbarray.Array(ids),
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		user_id = ANY(:user_id)`

	var usrs []dbUser
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &usrs); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUserSlice(usrs), nil
}

// QueryByEmail gets the specified user from the database by email.
func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email.Address,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		email = :email`

	var dbUsr dbUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUser(dbUsr), nil
}
