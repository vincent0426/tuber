package tripgrp

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/TSMC-Uber/server/business/core/trip"
	"github.com/TSMC-Uber/server/business/sys/validate"
	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (trip.QueryFilter, error) {
	values := r.URL.Query()

	var filter trip.QueryFilter

	if tripID := values.Get("trip_id"); tripID != "" {
		id, err := uuid.Parse(tripID)
		if err != nil {
			return trip.QueryFilter{}, validate.NewFieldsError("trip_id", err)
		}
		filter.WithTripID(id)
	}

	if driverID := values.Get("driver_id"); driverID != "" {
		id, err := uuid.Parse(driverID)
		if err != nil {
			return trip.QueryFilter{}, validate.NewFieldsError("driver_id", err)
		}
		filter.WithDriverID(id)
	}

	if passengerLimit := values.Get("passenger_limit"); passengerLimit != "" {
		pl, err := strconv.Atoi(passengerLimit)
		if err != nil {
			return trip.QueryFilter{}, validate.NewFieldsError("passenger_limit", err)
		}
		filter.WithPassengerLimit(pl)
	}

	if sourceID := values.Get("source_id"); sourceID != "" {
		id, err := uuid.Parse(sourceID)
		if err != nil {
			return trip.QueryFilter{}, validate.NewFieldsError("source_id", err)
		}
		filter.WithSourceID(id)
	}

	if destinationID := values.Get("destination_id"); destinationID != "" {
		id, err := uuid.Parse(destinationID)
		if err != nil {
			return trip.QueryFilter{}, validate.NewFieldsError("destination_id", err)
		}
		filter.WithDestinationID(id)
	}

	if start_start_date := values.Get("start_start_date"); start_start_date != "" {
		t, err := time.Parse(time.RFC3339, start_start_date)
		if err != nil {
			return trip.QueryFilter{}, validate.NewFieldsError("start_start_date", err)
		}
		filter.WithStartStartDate(t)
	}

	if end_start_date := values.Get("end_start_date"); end_start_date != "" {
		t, err := time.Parse(time.RFC3339, end_start_date)
		if err != nil {
			return trip.QueryFilter{}, validate.NewFieldsError("end_start_date", err)
		}
		filter.WithEndStartDate(t)
	}

	if err := filter.Validate(); err != nil {
		return trip.QueryFilter{}, err
	}

	return filter, nil
}

func parseFilterByUser(r *http.Request) (trip.QueryFilterByUser, error) {
	values := r.URL.Query()

	var filter trip.QueryFilterByUser

	if status := values.Get("status"); status != "" {
		// status should only have 3 values: "not_start", "in_trip", "finished"
		if status != "not_start" && status != "in_trip" && status != "finished" {
			return trip.QueryFilterByUser{}, errors.New("status should only have 3 values: not_start, in_trip, finished")
		}
		filter.WithStatus(status)
	}

	if isDriver := values.Get("is_driver"); isDriver != "" {
		isDriverBool, err := strconv.ParseBool(isDriver)
		if err != nil {
			return trip.QueryFilterByUser{}, validate.NewFieldsError("is_driver", err)
		}
		filter.IsDriver = &isDriverBool
	}

	if err := filter.Validate(); err != nil {
		return trip.QueryFilterByUser{}, err
	}

	return filter, nil
}
