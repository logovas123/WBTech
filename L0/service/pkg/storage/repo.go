package storage

import "service/pkg/model"

type OrderRepo interface {
	AddOrder(order *model.Order) (string, error)
	GetOrder(id string) (*model.Order, error)
	DeleteOrder(id string) error
}

type OrderRepoClose interface {
	OrderRepo
	Close()
}
