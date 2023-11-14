package tripdb

import (
	"github.com/Masterminds/squirrel"
	"github.com/TSMC-Uber/server/business/core/trip"
)

// ID             *uuid.UUID `validate:"omitempty"`
//
//	DriverID       *uuid.UUID `validate:"omitempty"`
//	PassengerLimit *int       `validate:"omitempty"`
//	SourceID       *uuid.UUID `validate:"omitempty"`
//	DestinationID  *uuid.UUID `validate:"omitempty"`
//	StartStartDate *time.Time `validate:"omitempty"`
//	EndStartDate   *time.Time `validate:"omitempty"`
func (s *Store) applyFilter(builder squirrel.SelectBuilder, filter trip.QueryFilter) squirrel.SelectBuilder {
	if filter.ID != nil {
		builder = builder.Where(squirrel.Eq{"id": *filter.ID})
	}

	if filter.DriverID != nil {
		builder = builder.Where(squirrel.Eq{"driver_id": *filter.DriverID})
	}

	if filter.PassengerLimit != nil {
		builder = builder.Where(squirrel.Eq{"passenger_limit": *filter.PassengerLimit})
	}

	if filter.SourceID != nil {
		builder = builder.Where(squirrel.Eq{"source_id": *filter.SourceID})
	}

	if filter.DestinationID != nil {
		builder = builder.Where(squirrel.Eq{"destination_id": *filter.DestinationID})
	}

	if filter.StartStartDate != nil {
		builder = builder.Where(squirrel.GtOrEq{"start_time": *filter.StartStartDate})
	}

	if filter.EndStartDate != nil {
		builder = builder.Where(squirrel.LtOrEq{"start_time": *filter.EndStartDate})
	}

	return builder
}
