package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// User represents information about an individual user.
type User struct {
	ID                 uuid.UUID
	Name               string
	Email              mail.Address
	ImageURL           string
	Bio                string
	AcceptNotification bool
	Sub                string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name               string
	Email              mail.Address
	ImageURL           string
	Bio                string
	AcceptNotification bool
	Sub                string
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Name               *string
	Email              *mail.Address
	ImageURL           *string
	Bio                *string
	AcceptNotification *bool
}
