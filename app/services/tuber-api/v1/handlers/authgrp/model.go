package authgrp

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/TSMC-Uber/server/business/core/auth"
	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/web/v1/mid"
)

type AppUser struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	ImageURL           string `json:"imageURL"`
	Bio                string `json:"bio"`
	AcceptNotification bool   `json:"acceptNotification"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

func toAppUser(usr user.User) AppUser {

	return AppUser{
		ID:                 usr.ID.String(),
		Name:               usr.Name,
		Email:              usr.Email.Address,
		ImageURL:           usr.ImageURL,
		Bio:                usr.Bio,
		AcceptNotification: usr.AcceptNotification,
		CreatedAt:          usr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          usr.UpdatedAt.Format(time.RFC3339),
	}
}

func toCoreNewUser(tokenInfo *auth.IDTokenInfo) (user.NewUser, error) {
	// upsert user and get user id here
	addr, err := mail.ParseAddress(tokenInfo.Email)
	if err != nil {
		return user.NewUser{}, mid.WrapError(fmt.Errorf("parsing email: %w", err))
	}
	usr := user.NewUser{
		Name:     tokenInfo.Name,
		Email:    *addr,
		ImageURL: tokenInfo.Picture,
		Sub:      tokenInfo.Sub,
	}

	return usr, nil
}
