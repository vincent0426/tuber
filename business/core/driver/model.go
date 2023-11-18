package driver

import (
	"time"

	"github.com/google/uuid"
)

type Driver struct {
	UserID    uuid.UUID `json:"user_id"`
	License   string    `json:"license"`
	Verified  bool      `json:"verified"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Color     string    `json:"color"`
	Plate     string    `json:"plate"`
	CreatedAt time.Time `json:"created_at"`
}

type NewDriver struct {
	UserID  uuid.UUID `json:"user_id"`
	License string    `json:"license"`
	Brand   string    `json:"brand"`
	Model   string    `json:"model"`
	Color   string    `json:"color"`
	Plate   string    `json:"plate"`
}

type FavoriteDriver struct {
	ID              uuid.UUID `json:"id"`
	DriverID        uuid.UUID `json:"driver_id"`
	DriverName      string    `json:"driver_name"`
	DriverImageURL  string    `json:"driver_image_url"`
	DriverBrand     string    `json:"driver_brand"`
	DriverModel     string    `json:"driver_model"`
	DriverColor     string    `json:"driver_color"`
	DriverPlate     string    `json:"driver_plate"`
	DriverCreatedAt time.Time `json:"driver_created_at"`
}
