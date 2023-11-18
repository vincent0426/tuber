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
