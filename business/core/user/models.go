package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// restrict language being one of the following, en, zh
type Language struct {
	Name string
}

// User represents information about an individual user.
type User struct {
	ID                 uuid.UUID
	Name               string
	Email              mail.Address
	Bio                string
	Lang               Language
	AcceptNotification bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name               string
	Email              mail.Address
	Bio                string
	Lang               Language
	AcceptNotification bool
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Name               *string
	Email              *mail.Address
	Bio                *string
	Lang               *Language
	AcceptNotification *bool
}
