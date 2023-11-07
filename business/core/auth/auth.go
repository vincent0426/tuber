// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package auth

import (
	"context"
	"crypto/sha256"

	"github.com/TSMC-Uber/server/business/core/user"

	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var ()

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Login(ctx context.Context, idToken string) (sessionToken string, err error)
	Logout(ctx context.Context, sessionToken string) error
	ValidateToken(ctx context.Context, tokenHash [32]byte) (user.User, error)
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

// Create inserts a new user into the database.
// func (c *Core) Login(ctx context.Context, idToken string) (sessionToken string, err error) {
// 	now := time.Now()

// 	token := Token{
// 		Hash: sha256.Sum256([]byte(idToken)),
// 	}

// 	if err := c.storer.Create(ctx, usr); err != nil {
// 		return User{}, fmt.Errorf("create: %w", err)
// 	}

// 	return usr, nil
// }

func (c *Core) ValidateToken(ctx context.Context, tokenPlaintext string) bool {
	// Calculate the SHA-256 hash of the plaintext token provided by the client. // Remember that this returns a byte *array* with length 32, not a slice. tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	// Set up the SQL query.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	user, err := c.storer.ValidateToken(ctx, tokenHash)
	if err != nil {
		return false
	}

	if user.ID != uuid.Nil {
		return true
	}

	return false
}
