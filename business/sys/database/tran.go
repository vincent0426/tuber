package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// dbBeginner implements the core coreTransaction interface,
type dbBeginner struct {
	sqlxDB *sqlx.DB
}

// NewBeginner constructs a value that implements the database
// beginner interface.
func NewBeginner(sqlxDB *sqlx.DB) Beginner {
	return &dbBeginner{
		sqlxDB: sqlxDB,
	}
}

// Begin start a transaction and returns a value that implements
// the core transactor interface.
func (db *dbBeginner) Begin() (Transaction, error) {
	return db.sqlxDB.Beginx()
}

// GetExtContext is a helper function that extracts the sqlx value
// from the core transactor interface for transactional use.
func GetExtContext(tx Transaction) (sqlx.ExtContext, error) {
	ec, ok := tx.(sqlx.ExtContext)
	if !ok {
		return nil, fmt.Errorf("Transactor(%T) not of a type *sql.Tx", tx)
	}

	return ec, nil
}

// Transaction represents a value that can commit or rollback a transaction.
type Transaction interface {
	Commit() error
	Rollback() error
}

// Beginner represents a value that can begin a transaction.
type Beginner interface {
	Begin() (Transaction, error)
}

// =============================================================================

type ctxKey int

const trKey ctxKey = 2

// Set stores a value that can manage a transaction.
func Set(ctx context.Context, tx Transaction) context.Context {
	return context.WithValue(ctx, trKey, tx)
}

// Get retrieves the value that can manage a transaction.
func Get(ctx context.Context) (Transaction, bool) {
	v, ok := ctx.Value(trKey).(Transaction)
	return v, ok
}

// =============================================================================

// ExecuteUnderTransaction is a helper function that can be used in tests and
// other apps to execute the core APIs under a transaction.
func ExecuteUnderTransaction(ctx context.Context, log *logger.Logger, bgn Beginner, fn func(tx Transaction) error) error {
	hasCommitted := false

	log.Info(ctx, "BEGIN TRANSACTION")
	tx, err := bgn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if !hasCommitted {
			log.Info(ctx, "ROLLBACK TRANSACTION")
		}

		if err := tx.Rollback(); err != nil {
			if errors.Is(err, sql.ErrTxDone) {
				return
			}
			log.Info(ctx, "ROLLBACK TRANSACTION", "ERROR", err)
		}
	}()

	if err := fn(tx); err != nil {
		return fmt.Errorf("EXECUTE TRANSACTION: %w", err)
	}

	log.Info(ctx, "COMMIT TRANSACTION")
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("COMMIT TRANSACTION: %w", err)
	}

	hasCommitted = true

	return nil
}
