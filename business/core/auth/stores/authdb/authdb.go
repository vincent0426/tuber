// Package authdb contains user related CRUD functionality.
package authdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/sys/database"
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

func (s *Store) Login(ctx context.Context, idToken string) (sessionToken string, err error) {
	return "", nil
}

func (s *Store) Logout(ctx context.Context, sessionToken string) error {
	return nil
}

func (s *Store) ValidateToken(ctx context.Context, tokenHash [32]byte) (user.User, error) {
	// user joins tokens
	sql, args, err := sq.
		Select("users.*").
		From("users").
		Join("tokens ON users.id = tokens.user_id").
		Where(sq.Eq{"tokens.hash": tokenHash[:]}).
		Where(sq.Gt{"tokens.expiry": time.Now()}).
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
