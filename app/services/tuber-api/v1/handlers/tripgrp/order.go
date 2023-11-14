package tripgrp

import (
	"errors"
	"net/http"

	"github.com/TSMC-Uber/server/business/core/trip"
	"github.com/TSMC-Uber/server/business/data/order"
	"github.com/TSMC-Uber/server/business/sys/validate"
)

var orderByFields = map[string]struct{}{
	trip.OrderByStatus:    {},
	trip.OrderByStartTime: {},
	trip.OrderByCreatedAt: {},
}

func parseOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, trip.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}
