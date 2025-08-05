package application

import "github.com/Babushkin05/wb-orders-service/internal/domain/model"

type OrdersRepository interface {
	Get(orderUID string) (string, error)
	Store(model model.Order) error
}
