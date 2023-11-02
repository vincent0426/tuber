package userdb

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/TSMC-Uber/server/business/core/user"
)

func (s *Store) applyFilter(builder squirrel.SelectBuilder, filter user.QueryFilter) squirrel.SelectBuilder {
	if filter.ID != nil {
		builder = builder.Where(squirrel.Eq{"user_id": *filter.ID})
	}

	if filter.Name != nil {
		builder = builder.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *filter.Name))
	}

	if filter.Email != nil {
		builder = builder.Where(squirrel.Eq{"email": (*filter.Email).String()})
	}

	if filter.StartCreatedDate != nil {
		builder = builder.Where(squirrel.GtOrEq{"date_created": *filter.StartCreatedDate})
	}

	if filter.EndCreatedDate != nil {
		builder = builder.Where(squirrel.LtOrEq{"date_created": *filter.EndCreatedDate})
	}

	return builder
}
