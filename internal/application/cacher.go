package application

import "github.com/Babushkin05/wb-orders-service/internal/domain/model"

type Casher interface {
	Cache(order model.Order) error
	GetOrderFromCache(orderUID string) (*model.Order, error)
}
