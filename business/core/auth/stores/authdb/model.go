package authdb

import (
	"net/mail"
	"time"

	"github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/google/uuid"
)

type dbSessionToken struct {
	Hash   []byte    `db:"hash"`
	UserID uuid.UUID `db:"user_id"`
	Expiry time.Time `db:"expiry"`
	Scope  string    `db:"scope"`
}

func toDBSessionToken(app auth.SessionToken) dbSessionToken {
	return dbSessionToken{
		Hash:   app.Hash,
		UserID: app.UserID,
		Expiry: app.Expiry,
		Scope:  app.Scope,
	}
}

// should not be like this
type dbUser struct {
	ID                 uuid.UUID `db:"id"`
	Name               string    `db:"name"`
	Email              string    `db:"email"`
	Bio                string    `db:"bio"`
	Lang               string    `db:"language"`
	AcceptNotification bool      `db:"accept_notification"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

func toCoreUser(dbUsr dbUser) user.User {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	usr := user.User{
		ID:                 dbUsr.ID,
		Name:               dbUsr.Name,
		Email:              addr,
		Bio:                dbUsr.Bio,
		Lang:               user.Language{Name: dbUsr.Lang},
		AcceptNotification: dbUsr.AcceptNotification,
		CreatedAt:          dbUsr.CreatedAt.In(time.Local),
		UpdatedAt:          dbUsr.UpdatedAt.In(time.Local),
	}

	return usr
}
