package application

import (
	"errors"

	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
)

type OrdersService interface {
	GetOrder(orderUID string) (model.Order, error)
	SaveOrder(order *model.Order) error
}

type ordersService struct {
	cacher           Cacher
	ordersRepository OrdersRepository
}

var _ OrdersService = &ordersService{}

func NewOrdersService(casher Cacher, ordersRepository OrdersRepository) OrdersService {
	return &ordersService{
		cacher:           casher,
		ordersRepository: ordersRepository,
	}
}

func (s *ordersService) GetOrder(orderUID string) (model.Order, error) {
	if orderUID == "" {
		return model.Order{}, errors.New("orderUID is empty")
	}

	order, err := s.cacher.GetOrderFromCache(orderUID)
	if err != nil {
		order, err = s.ordersRepository.Get(orderUID)
		if err != nil {
			return model.Order{}, err
		}
	}

	return order, nil
}

func (s *ordersService) SaveOrder(order *model.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}

	err := s.cacher.Cache(order)
	if err != nil {
		return err
	}

	return s.ordersRepository.Store(order)
}
