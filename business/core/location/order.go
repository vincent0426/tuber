package location

import "github.com/TSMC-Uber/server/business/data/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByName, order.ASC)

const (
	OrderByName = "name"
)
