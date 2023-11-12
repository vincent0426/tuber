package usergrp

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/sys/validate"
	"github.com/TSMC-Uber/server/business/web/v1/mid"
)

// AppUser represents information about an individual user.
type AppUser struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
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
		Bio:                usr.Bio,
		AcceptNotification: usr.AcceptNotification,
		CreatedAt:          usr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          usr.UpdatedAt.Format(time.RFC3339),
	}
}

// =============================================================================

// AppNewUser contains information needed to create a new user.
type AppNewUser struct {
	Name               string `json:"name" binding:"required"`
	Email              string `json:"email" binding:"required,email"`
	Bio                string `json:"bio"`
	AcceptNotification bool   `json:"acceptNotification"`
}

func toCoreNewUser(app AppNewUser) (user.NewUser, error) {
	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return user.NewUser{}, mid.WrapError(fmt.Errorf("parsing email: %w", err))
	}

	usr := user.NewUser{
		Name:               app.Name,
		Email:              *addr,
		Bio:                app.Bio,
		AcceptNotification: app.AcceptNotification,
	}

	return usr, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewUser) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

// =============================================================================

// AppUpdateUser contains information needed to update a user.
type AppUpdateUser struct {
	Name               *string `json:"name"`
	Email              *string `json:"email" validate:"omitempty,email"`
	Bio                *string `json:"bio"`
	AcceptNotification *bool   `json:"acceptNotification"`
}

func toCoreUpdateUser(app AppUpdateUser) (user.UpdateUser, error) {
	var addr *mail.Address
	if app.Email != nil {
		var err error
		addr, err = mail.ParseAddress(*app.Email)
		if err != nil {
			return user.UpdateUser{}, mid.WrapError(fmt.Errorf("parsing email: %w", err))
		}
	}

	nu := user.UpdateUser{
		Name:               app.Name,
		Email:              addr,
		Bio:                app.Bio,
		AcceptNotification: app.AcceptNotification,
	}

	return nu, nil
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateUser) Validate() error {
	if err := validate.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// =============================================================================

// AppSummary represents information about an individual user and their products.
type AppSummary struct {
	UserID     string  `json:"userID"`
	UserName   string  `json:"userName"`
	TotalCount int     `json:"totalCount"`
	TotalCost  float64 `json:"totalCost"`
}
