package drivergrp

import (
	"errors"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/driver"
	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/TSMC-Uber/server/business/sys/validate"
)

var orderByFields = map[string]struct{}{
	driver.OrderByBrand: {},
	driver.OrderByModel: {},
	driver.OrderByColor: {},
}

var orderByFieldsFavoriteDriver = map[string]struct{}{
	driver.OrderByBrandFavoriteDriver: {},
	driver.OrderByModelFavoriteDriver: {},
	driver.OrderByColorFavoriteDriver: {},
}

func parseOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, driver.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}

func parsefavoriteDriverOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, driver.DefaultOrderByFavoriteDriver)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFieldsFavoriteDriver[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}
