package storage

import "service/pkg/model"

// данный интерфейс определяет поведение БД, тем самым к сервису может быть подключена любая бд (например данным интерфейсом могут быть описаны бд и кеш)
type OrderRepo interface {
	AddOrder(order *model.Order) (string, error)
	GetOrder(id string) (*model.Order, error)
	DeleteOrder(id string) error
}

// данный интерфейс относится к БД, т.к. пул соединений к базе необходимо закрывать (я сам решил добавить такой интерфейс для удобства)
type OrderRepoClose interface {
	OrderRepo
	Close()
}
