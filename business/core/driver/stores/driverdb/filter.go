package driverdb

import (
	"github.com/Masterminds/squirrel"
	"github.com/TSMC-Uber/server/business/core/driver"
)

func (s *Store) applyFilter(builder squirrel.SelectBuilder, filter driver.QueryFilter) squirrel.SelectBuilder {
	if filter.DriverID != nil {
		builder = builder.Where(squirrel.Eq{"driver.user_id": *filter.DriverID})
	}

	if filter.Brand != nil {
		builder = builder.Where(squirrel.Eq{"driver.brand": *filter.Brand})
	}

	if filter.Model != nil {
		builder = builder.Where(squirrel.Eq{"driver.model": *filter.Model})
	}

	if filter.Color != nil {
		builder = builder.Where(squirrel.Eq{"driver.color": *filter.Color})
	}

	return builder
}
