package drivergrp

import (
	"time"

	"github.com/TSMC-Uber/server/business/core/driver"
)

// AppDriver represents information about an individual location.
type AppDriver struct {
	UserID    string `json:"user_id"`
	License   string `json:"license"`
	Verified  bool   `json:"verified"`
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	Color     string `json:"color"`
	Plate     string `json:"plate"`
	CreatedAt string `json:"created_at"`
}

func toAppDriver(driver driver.Driver) AppDriver {

	return AppDriver{
		UserID:    driver.UserID.String(),
		License:   driver.License,
		Verified:  driver.Verified,
		Brand:     driver.Brand,
		Model:     driver.Model,
		Color:     driver.Color,
		Plate:     driver.Plate,
		CreatedAt: driver.CreatedAt.Format(time.RFC3339),
	}
}

// =============================================================================

// AppNewUser contains information needed to create a new user.
type AppNewDriver struct {
	License string `json:"license" binding:"required"`
	Brand   string `json:"brand" binding:"required"`
	Model   string `json:"model" binding:"required"`
	Color   string `json:"color" binding:"required"`
	Plate   string `json:"plate" binding:"required"`
}

func toCoreNewDriver(app AppNewDriver) (driver.NewDriver, error) {

	driver := driver.NewDriver{
		License: app.License,
		Brand:   app.Brand,
		Model:   app.Model,
		Color:   app.Color,
		Plate:   app.Plate,
	}

	return driver, nil
}

type AppFavoriteDriver struct {
	ID              string `json:"id"`
	DriverID        string `json:"driver_id"`
	DriverName      string `json:"driver_name"`
	DriverImageURL  string `json:"driver_image_url"`
	DriverBrand     string `json:"driver_brand"`
	DriverModel     string `json:"driver_model"`
	DriverColor     string `json:"driver_color"`
	DriverPlate     string `json:"driver_plate"`
	DriverCreatedAt string `json:"driver_created_at"`
}

func toAppFavoriteDriver(driver driver.FavoriteDriver) AppFavoriteDriver {

	return AppFavoriteDriver{
		ID:              driver.ID.String(),
		DriverID:        driver.DriverID.String(),
		DriverName:      driver.DriverName,
		DriverImageURL:  driver.DriverImageURL,
		DriverBrand:     driver.DriverBrand,
		DriverModel:     driver.DriverModel,
		DriverColor:     driver.DriverColor,
		DriverPlate:     driver.DriverPlate,
		DriverCreatedAt: driver.DriverCreatedAt.Format(time.RFC3339),
	}
}
