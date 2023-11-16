package locationdb

import (
	"github.com/Masterminds/squirrel"
	"github.com/TSMC-Uber/server/business/core/location"
)

// ID             *uuid.UUID `validate:"omitempty"`
//
//	DriverID       *uuid.UUID `validate:"omitempty"`
//	PassengerLimit *int       `validate:"omitempty"`
//	SourceID       *uuid.UUID `validate:"omitempty"`
//	DestinationID  *uuid.UUID `validate:"omitempty"`
//	StartStartDate *time.Time `validate:"omitempty"`
//	EndStartDate   *time.Time `validate:"omitempty"`
func (s *Store) applyFilter(builder squirrel.SelectBuilder, filter location.QueryFilter) squirrel.SelectBuilder {
	if filter.ID != nil {
		builder = builder.Where(squirrel.Eq{"id": *filter.ID})
	}

	return builder
}
