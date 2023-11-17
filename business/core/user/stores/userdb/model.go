package userdb

import (
	"net/mail"
	"time"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/google/uuid"
)

// dbUser represent the structure we need for moving data
// between the app and the database.
type dbUser struct {
	ID                 uuid.UUID `db:"id"`
	Name               string    `db:"name"`
	Email              string    `db:"email"`
	ImageURL           string    `db:"image_url"`
	Bio                string    `db:"bio"`
	AcceptNotification bool      `db:"accept_notification"`
	Sub                string    `db:"sub"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

func toDBUser(usr user.User) dbUser {
	return dbUser{
		ID:                 usr.ID,
		Name:               usr.Name,
		Email:              usr.Email.Address,
		ImageURL:           usr.ImageURL,
		Bio:                usr.Bio,
		AcceptNotification: usr.AcceptNotification,
		Sub:                usr.Sub,
		CreatedAt:          usr.CreatedAt.UTC(),
		UpdatedAt:          usr.UpdatedAt.UTC(),
	}
}

func toCoreUser(dbUsr dbUser) user.User {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	usr := user.User{
		ID:                 dbUsr.ID,
		Name:               dbUsr.Name,
		Email:              addr,
		ImageURL:           dbUsr.ImageURL,
		Bio:                dbUsr.Bio,
		AcceptNotification: dbUsr.AcceptNotification,
		CreatedAt:          dbUsr.CreatedAt.In(time.Local),
		UpdatedAt:          dbUsr.UpdatedAt.In(time.Local),
	}

	return usr
}

func toCoreUserSlice(dbUsers []dbUser) []user.User {
	usrs := make([]user.User, len(dbUsers))
	for i, dbUsr := range dbUsers {
		usrs[i] = toCoreUser(dbUsr)
	}
	return usrs
}
