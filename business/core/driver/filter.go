package driver

import (
	"fmt"

	"github.com/TSMC-Uber/server/business/sys/validate"
	"github.com/google/uuid"
)

type QueryFilter struct {
	DriverID *uuid.UUID `validate:"omitempty"`
	Brand    *string    `validate:"omitempty"`
	Model    *string    `validate:"omitempty"`
	Color    *string    `validate:"omitempty"`
}

// Validate checks the data in the model is considered clean.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithDriverID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithDriverID(driverID uuid.UUID) {
	qf.DriverID = &driverID
}

// WithBrand sets the Name field of the QueryFilter value.
func (qf *QueryFilter) WithBrand(brand string) {
	qf.Brand = &brand
}

// WithModel sets the Email field of the QueryFilter value.
func (qf *QueryFilter) WithModel(model string) {
	qf.Model = &model
}

// WithColor sets the DateCreated field of the QueryFilter value.
func (qf *QueryFilter) WithColor(color string) {
	qf.Color = &color
}

type QueryFilterFavoriteDriver struct {
	UserID   *uuid.UUID `validate:"omitempty"`
	DriverID *uuid.UUID `validate:"omitempty"`
	Brand    *string    `validate:"omitempty"`
	Model    *string    `validate:"omitempty"`
	Color    *string    `validate:"omitempty"`
}

// Validate checks the data in the model is considered clean.
func (qf *QueryFilterFavoriteDriver) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func (qf *QueryFilterFavoriteDriver) WithUserID(userID uuid.UUID) {
	qf.UserID = &userID
}

// WithDriverID sets the ID field of the QueryFilter value.
func (qf *QueryFilterFavoriteDriver) WithDriverID(driverID uuid.UUID) {
	qf.DriverID = &driverID
}

// WithBrand sets the Name field of the QueryFilter value.
func (qf *QueryFilterFavoriteDriver) WithBrand(brand string) {
	qf.Brand = &brand
}

// WithModel sets the Email field of the QueryFilter value.
func (qf *QueryFilterFavoriteDriver) WithModel(model string) {
	qf.Model = &model
}

// WithColor sets the DateCreated field of the QueryFilter value.
func (qf *QueryFilterFavoriteDriver) WithColor(color string) {
	qf.Color = &color
}
