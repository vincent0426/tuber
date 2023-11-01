package userdb

import (
	"fmt"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/TSMC-Uber/server/business/data/order"
)

var orderByFields = map[string]string{
	user.OrderByID:    "id",
	user.OrderByName:  "name",
	user.OrderByEmail: "email",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
