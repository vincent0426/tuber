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

type dbFavoriteDriver struct {
	FavoriteDriverID uuid.UUID `db:"id"`
	DriverID         uuid.UUID `db:"driver_id"`
	DriverName       string    `db:"driver_name"`
	DriverImageURL   string    `db:"driver_image_url"`
	DriverBrand      string    `db:"driver_brand"`
	DriverModel      string    `db:"driver_model"`
	DriverColor      string    `db:"driver_color"`
	DriverPlate      string    `db:"driver_plate"`
	DriverCreatedAt  time.Time `db:"driver_created_at"`
}

func toCoreFavoriteDriver(dbDriver dbFavoriteDriver) driver.FavoriteDriver {
	driver := driver.FavoriteDriver{
		ID:              dbDriver.DriverID,
		DriverID:        dbDriver.DriverID,
		DriverName:      dbDriver.DriverName,
		DriverImageURL:  dbDriver.DriverImageURL,
		DriverBrand:     dbDriver.DriverBrand,
		DriverModel:     dbDriver.DriverModel,
		DriverColor:     dbDriver.DriverColor,
		DriverPlate:     dbDriver.DriverPlate,
		DriverCreatedAt: dbDriver.DriverCreatedAt.In(time.Local),
	}

	return driver
}

func toCoreFavoriteDriverSlice(dbFavoriteDrivers []dbFavoriteDriver) []driver.FavoriteDriver {
	drivers := make([]driver.FavoriteDriver, len(dbFavoriteDrivers))
	for i, dbDriver := range dbFavoriteDrivers {
		drivers[i] = toCoreFavoriteDriver(dbDriver)
	}
	return drivers
}
