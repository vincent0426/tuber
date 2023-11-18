package driverdb

import (
	"fmt"

	"github.com/TSMC-Uber/server/business/core/driver"
	"github.com/TSMC-Uber/server/business/data/order"
)

var orderByFields = map[string]string{
	driver.OrderByBrand: "brand",
	driver.OrderByModel: "model",
	driver.OrderByColor: "color",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return fmt.Sprintf("%s %s", by, orderBy.Direction), nil
}

var orderByFieldsFavoriteDriver = map[string]string{
	driver.OrderByBrandFavoriteDriver: "driver_brand",
	driver.OrderByModelFavoriteDriver: "driver_model",
	driver.OrderByColorFavoriteDriver: "driver_color",
}

func orderByClauseFavoriteDriver(orderBy order.By) (string, error) {
	by, exists := orderByFieldsFavoriteDriver[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return fmt.Sprintf("%s %s", by, orderBy.Direction), nil
}
