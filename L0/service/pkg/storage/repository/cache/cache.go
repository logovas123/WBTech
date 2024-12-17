package cache

import (
	"log/slog"
	"sync"

	"service/pkg/model"
	"service/pkg/storage"
)

type OrderCacheRepository struct {
	Orders map[string]*model.Order
	Mu     *sync.RWMutex
}

func NewCashe() *OrderCacheRepository {
	return &OrderCacheRepository{
		Orders: map[string]*model.Order{},
		Mu:     &sync.RWMutex{},
	}
}

func (cacheOrders *OrderCacheRepository) AddOrder(order *model.Order) (string, error) {
	cacheOrders.Mu.RLock()
	if _, ok := cacheOrders.Orders[order.OrderUID]; ok {
		cacheOrders.Mu.RUnlock()
		slog.Error("error exist", "id", order.OrderUID)
		return "", storage.ErrorOrderExist
	}
	cacheOrders.Mu.RUnlock()

	cacheOrders.Mu.Lock()
	cacheOrders.Orders[order.OrderUID] = order
	cacheOrders.Mu.Unlock()

	slog.Info("order add in cache success")

	return order.OrderUID, nil
}

func (cacheOrders *OrderCacheRepository) GetOrder(id string) (*model.Order, error) {
	cacheOrders.Mu.RLock()
	if _, ok := cacheOrders.Orders[id]; !ok {
		cacheOrders.Mu.RUnlock()
		slog.Info("error get order:",
			"id", id,
			"error", storage.ErrorOrderNotExist,
		)
		return nil, storage.ErrorOrderNotExist
	}

	order := cacheOrders.Orders[id]
	cacheOrders.Mu.RUnlock()

	return order, nil
}

func (cacheOrders *OrderCacheRepository) DeleteOrder(id string) error {
	cacheOrders.Mu.RLock()
	if _, ok := cacheOrders.Orders[id]; !ok {
		cacheOrders.Mu.RUnlock()
		return storage.ErrorOrderNotExist
	}
	cacheOrders.Mu.RUnlock()

	cacheOrders.Mu.Lock()
	delete(cacheOrders.Orders, id)
	cacheOrders.Mu.Unlock()

	return nil
}
