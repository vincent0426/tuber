package driver

import "github.com/TSMC-Uber/server/business/data/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByBrand, order.ASC)

const (
	OrderByBrand = "brand"
	OrderByModel = "model"
	OrderByColor = "color"
)

var DefaultOrderByFavoriteDriver = order.NewBy(OrderByBrandFavoriteDriver, order.ASC)

const (
	OrderByBrandFavoriteDriver = "driver_brand"
	OrderByModelFavoriteDriver = "driver_model"
	OrderByColorFavoriteDriver = "driver_color"
)
