package driverdb

import (
	"time"

	"github.com/TSMC-Uber/server/business/core/driver"
	"github.com/google/uuid"
)

type dbDriver struct {
	UserID    uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	ImageURL  string    `db:"image_url"`
	License   string    `db:"license"`
	Verified  bool      `db:"verified"`
	Brand     string    `db:"brand"`
	Model     string    `db:"model"`
	Color     string    `db:"color"`
	Plate     string    `db:"plate"`
	CreatedAt time.Time `db:"driver_created_at"`
}

func toDBDriver(driver driver.Driver) dbDriver {
	return dbDriver{
		UserID:    driver.UserID,
		License:   driver.License,
		Verified:  driver.Verified,
		Brand:     driver.Brand,
		Model:     driver.Model,
		Color:     driver.Color,
		Plate:     driver.Plate,
		CreatedAt: driver.CreatedAt.UTC(),
	}
}

func toCoreDriver(dbDriver dbDriver) driver.Driver {

	trip := driver.Driver{
		UserID:    dbDriver.UserID,
		License:   dbDriver.License,
		Verified:  dbDriver.Verified,
		Brand:     dbDriver.Brand,
		Model:     dbDriver.Model,
		Color:     dbDriver.Color,
		Plate:     dbDriver.Plate,
		CreatedAt: dbDriver.CreatedAt.In(time.Local),
	}

	return trip
}

func toCoreDriverSlice(dbDrivers []dbDriver) []driver.Driver {
	drivers := make([]driver.Driver, len(dbDrivers))
	for i, dbDriver := range dbDrivers {
		drivers[i] = toCoreDriver(dbDriver)
	}
	return drivers
}
