package cache

import (
	"log/slog"
	"sync"

	"service/pkg/model"
	"service/pkg/storage"
)

// здесь применяется паттерн репозиторий

/*
структура кеша, использую мапу и мьютекс для избежания конкурентного доступа к мапе.
мапа хранит Order, который представляет собой структуру, описанную в папке model.
для кеша реализованы основные методы, позволяющие получить, добавить и удалить заказ (используются не все методы)
*/

type OrderCacheRepository struct {
	Orders map[string]*model.Order
	Mu     *sync.RWMutex
}

// создаём кеш
func NewCashe() *OrderCacheRepository {
	return &OrderCacheRepository{
		Orders: map[string]*model.Order{},
		Mu:     &sync.RWMutex{},
	}
}

// метод добавления заказа
func (cacheOrders *OrderCacheRepository) AddOrder(order *model.Order) (string, error) {
	cacheOrders.Mu.RLock()
	// проверяем существовал ли заказ, который мы хотим добавить
	if _, ok := cacheOrders.Orders[order.OrderUID]; ok {
		cacheOrders.Mu.RUnlock()
		slog.Error("error exist", "id", order.OrderUID)
		return "", storage.ErrorOrderExist // если заказ существовал возвращаем кастомную ошибку
	}
	cacheOrders.Mu.RUnlock()

	cacheOrders.Mu.Lock()
	cacheOrders.Orders[order.OrderUID] = order // добалвяем заказ
	cacheOrders.Mu.Unlock()

	slog.Info("order add in cache success")

	return order.OrderUID, nil
}

// метод получения заказа
func (cacheOrders *OrderCacheRepository) GetOrder(id string) (*model.Order, error) {
	cacheOrders.Mu.RLock()
	// полуаем заказ
	if _, ok := cacheOrders.Orders[id]; !ok {
		cacheOrders.Mu.RUnlock()
		slog.Info("error get order:",
			"id", id,
			"error", storage.ErrorOrderNotExist,
		)
		return nil, storage.ErrorOrderNotExist // если заказ не существует, возвращаем ошибку
	}

	order := cacheOrders.Orders[id] // получаем заказ
	cacheOrders.Mu.RUnlock()

	return order, nil // возвращаем полученный заказ
}

// метод для удаления заказа
func (cacheOrders *OrderCacheRepository) DeleteOrder(id string) error {
	cacheOrders.Mu.RLock()
	// проверяем существование заказа
	if _, ok := cacheOrders.Orders[id]; !ok {
		cacheOrders.Mu.RUnlock()
		return storage.ErrorOrderNotExist
	}
	cacheOrders.Mu.RUnlock()

	cacheOrders.Mu.Lock()
	delete(cacheOrders.Orders, id) // удаляем заказ
	cacheOrders.Mu.Unlock()

	return nil
}
