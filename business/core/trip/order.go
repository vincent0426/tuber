package trip

import "github.com/TSMC-Uber/server/business/data/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByStartTime, order.ASC)

const (
	OrderByStatus    = "status"
	OrderByStartTime = "start_time"
	OrderByCreatedAt = "created_at"
)
