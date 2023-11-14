// Package authdb contains user related CRUD functionality.
package authdb

import (
	"context"

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

func (s *Store) Login(ctx context.Context, idToken string) (sessionToken string, err error) {
	return "", nil
}

func (s *Store) Logout(ctx context.Context, sessionToken string) error {
	return nil
}
