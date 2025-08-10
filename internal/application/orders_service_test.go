package application

import (
	"context"
	"errors"
	"testing"

	"github.com/Babushkin05/wb-orders-service/internal/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ----- Моки -----
type mockCacher struct{ mock.Mock }

func (m *mockCacher) Cache(order *model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *mockCacher) GetOrderFromCache(orderUID string) (model.Order, error) {
	args := m.Called(orderUID)
	return args.Get(0).(model.Order), args.Error(1)
}
func (m *mockCacher) WarmUp() error {
	args := m.Called()
	return args.Error(0)
}

type mockOrdersRepository struct{ mock.Mock }

func (m *mockOrdersRepository) Get(orderUID string) (model.Order, error) {
	args := m.Called(orderUID)
	return args.Get(0).(model.Order), args.Error(1)
}
func (m *mockOrdersRepository) Store(order *model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *mockOrdersRepository) SaveInboxMessage(_ context.Context, _, _, _ string) error {
	return nil
}
func (m *mockOrdersRepository) FetchUnprocessedInboxMessages(_ context.Context, _ int) ([]model.InboxMessage, error) {
	return nil, nil
}
func (m *mockOrdersRepository) MarkInboxMessageProcessed(_ context.Context, _ string) error {
	return nil
}

// ----- Тесты -----
func TestGetOrder_FromCacheSuccess(t *testing.T) {
	cacher := new(mockCacher)
	repo := new(mockOrdersRepository)

	uid := uuid.New()
	expectedOrder := model.Order{OrderUID: uid}
	cacher.On("GetOrderFromCache", uid.String()).Return(expectedOrder, nil)

	service := NewOrdersService(cacher, repo)

	order, err := service.GetOrder(uid.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestGetOrder_FromRepoAfterCacheMiss(t *testing.T) {
	cacher := new(mockCacher)
	repo := new(mockOrdersRepository)

	uid := uuid.New()
	expectedOrder := model.Order{OrderUID: uid}
	cacher.On("GetOrderFromCache", uid.String()).Return(model.Order{}, errors.New("not found"))
	repo.On("Get", uid.String()).Return(expectedOrder, nil)

	service := NewOrdersService(cacher, repo)

	order, err := service.GetOrder(uid.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestGetOrder_EmptyUID(t *testing.T) {
	service := NewOrdersService(nil, nil)
	order, err := service.GetOrder("")

	assert.Error(t, err)
	assert.Empty(t, order)
}

func TestGetOrder_RepoError(t *testing.T) {
	cacher := new(mockCacher)
	repo := new(mockOrdersRepository)

	uid := uuid.New()
	cacher.On("GetOrderFromCache", uid.String()).Return(model.Order{}, errors.New("not found"))
	repo.On("Get", uid.String()).Return(model.Order{}, errors.New("db error"))

	service := NewOrdersService(cacher, repo)

	order, err := service.GetOrder(uid.String())

	assert.Error(t, err)
	assert.Empty(t, order.OrderUID)
}

func TestSaveOrder_Success(t *testing.T) {
	cacher := new(mockCacher)
	repo := new(mockOrdersRepository)

	uid := uuid.New()
	order := &model.Order{OrderUID: uid}
	cacher.On("Cache", order).Return(nil)
	repo.On("Store", order).Return(nil)

	service := NewOrdersService(cacher, repo)

	err := service.SaveOrder(order)

	assert.NoError(t, err)
}

func TestSaveOrder_NilOrder(t *testing.T) {
	service := NewOrdersService(nil, nil)
	err := service.SaveOrder(nil)

	assert.Error(t, err)
}

func TestSaveOrder_CacheError(t *testing.T) {
	cacher := new(mockCacher)
	repo := new(mockOrdersRepository)

	uid := uuid.New()
	order := &model.Order{OrderUID: uid}
	cacher.On("Cache", order).Return(errors.New("cache error"))

	service := NewOrdersService(cacher, repo)

	err := service.SaveOrder(order)

	assert.Error(t, err)
}

func TestSaveOrder_RepoError(t *testing.T) {
	cacher := new(mockCacher)
	repo := new(mockOrdersRepository)

	uid := uuid.New()
	order := &model.Order{OrderUID: uid}
	cacher.On("Cache", order).Return(nil)
	repo.On("Store", order).Return(errors.New("db error"))

	service := NewOrdersService(cacher, repo)

	err := service.SaveOrder(order)

	assert.Error(t, err)
}
