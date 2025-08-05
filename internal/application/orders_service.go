package application

import (
	"errors"

	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
)

type OrdersService interface {
	GetOrder(orderUID string) (model.Order, error)
	SaveOrder(order model.Order) error
}

type Service struct {
	Casher           *Casher
	OrdersRepository *OrdersRepository
}

func NewService(casher *Casher, ordersRepository *OrdersRepository) *Service {
	return &Service{
		Casher:           casher,
		OrdersRepository: ordersRepository,
	}
}

_, ok = s.OrdersRepository.(*OrdersRepository)

func (s *Service) GetOrder(orderUID string) (model.Order, error) {
	if orderUID == "" {
		return model.Order{}, errors.New("orderUID is empty")
	}

	order, err := s.Cacher.Get(orderUID)
	if(err != nil) {
		logger.Log.Error(err)
		return nil, err
	}

	if(order == nil) {
		order, err = s.OrdersRepository.Get(orderUID)
		if(err != nil) {
			logger.Log.Error(err)
			return nil, err
		}
	}
	return order, nil
}

func (s *Service) SaveOrder(order model.Order) error {
	if(order == nil) {
		return errors.New("order is nil")
	}
	if(order.UID == "") {
		return errors.New("orderUID is empty")
	}

	return s.OrdersRepository.Store(order)
}