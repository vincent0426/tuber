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
	Bio                string    `db:"bio"`
	Lang               string    `db:"language"`
	AcceptNotification bool      `db:"accept_notification"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

func toDBUser(usr user.User) dbUser {
	return dbUser{
		ID:                 usr.ID,
		Name:               usr.Name,
		Email:              usr.Email.Address,
		Bio:                usr.Bio,
		Lang:               usr.Lang.Name,
		AcceptNotification: usr.AcceptNotification,
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
		Bio:                dbUsr.Bio,
		Lang:               user.Language{Name: dbUsr.Lang},
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
