package auth

import (
	"time"

	"github.com/google/uuid"
)

// User represents information about an individual user.
type Token struct {
	Hash   []byte
	UserID uuid.UUID
	Expiry time.Time
	Scope  string
}

// NewUser contains information needed to create a new user.
type NewToken struct {
	Hash   []byte
	UserID uuid.UUID
	Expiry time.Time
	Scope  string
}
