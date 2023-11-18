package drivergrp

import (
	"net/http"

	"github.com/TSMC-Uber/server/business/core/driver"
	"github.com/TSMC-Uber/server/business/sys/validate"
	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (driver.QueryFilter, error) {
	values := r.URL.Query()

	var filter driver.QueryFilter

	if driverID := values.Get("driver_id"); driverID != "" {
		id, err := uuid.Parse(driverID)
		if err != nil {
			return driver.QueryFilter{}, validate.NewFieldsError("driver_id", err)
		}
		filter.WithDriverID(id)
	}

	if brand := values.Get("brand"); brand != "" {
		filter.WithBrand(brand)
	}

	if model := values.Get("model"); model != "" {
		filter.WithModel(model)
	}

	if color := values.Get("color"); color != "" {
		filter.WithColor(color)
	}

	if err := filter.Validate(); err != nil {
		return driver.QueryFilter{}, err
	}

	return filter, nil
}
