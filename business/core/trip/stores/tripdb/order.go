package tripdb

import (
	"fmt"

	"github.com/TSMC-Uber/server/business/core/trip"
	"github.com/TSMC-Uber/server/business/data/order"
)

var orderByFields = map[string]string{
	trip.OrderByStatus:    "status",
	trip.OrderByStartTime: "start_time",
	trip.OrderByCreatedAt: "created_at",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		fmt.Println("orderByFields", orderByFields)
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return fmt.Sprintf("%s %s", by, orderBy.Direction), nil
}
