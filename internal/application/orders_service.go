package application

import (
	"errors"

	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
)

type OrdersService interface {
	GetOrder(orderUID string) (*model.Order, error)
	SaveOrder(order *model.Order) error
}

type Service struct {
	cacher           Cacher
	ordersRepository OrdersRepository
}

var _ OrdersService = &Service{}

func NewService(casher Cacher, ordersRepository OrdersRepository) *Service {
	return &Service{
		cacher:           casher,
		ordersRepository: ordersRepository,
	}
}

func (s *Service) GetOrder(orderUID string) (*model.Order, error) {
	if orderUID == "" {
		return &model.Order{}, errors.New("orderUID is empty")
	}

	order, err := s.cacher.GetOrderFromCache(orderUID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		order, err = s.ordersRepository.Get(orderUID)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}

func (s *Service) SaveOrder(order *model.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}
	if order.OrderUID == "" {
		return errors.New("orderUID is empty")
	}

	err := s.cacher.Cache(order)
	if err != nil {
		return err
	}

	return s.ordersRepository.Store(order)
}
