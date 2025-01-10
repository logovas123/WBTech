package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"service/pkg/model"
	"service/pkg/storage"
)

// метод для обработки сообщения
func (c Consumer) handleMessage(ctx context.Context, message []byte) error {
	order := &model.Order{} // создаём пустую структуру заказа (описание структуры Order находится в папке model)

	err := json.Unmarshal(message, order) // анмаршалим полученное сообщение в структуру
	if err != nil {
		slog.Error("error unmarshal message:", "error", err)
		return err
	}
	slog.Info("msg unmarshal success", "order id", order.OrderUID)

	err = c.addToCache(ctx, order) // добавляем заказ в кеш
	if err != nil {
		slog.Error("error add to cache")
		return err
	}
	slog.Info("order add in cache success")

	err = c.addToDB(ctx, order) // добавляем заказ в бд
	if err != nil {
		slog.Error("error add to db", "error", err)
		return err
	}
	slog.Info("order add in db success")

	return nil
}

// метод добавления записи в кеш
func (c Consumer) addToCache(ctx context.Context, order *model.Order) error {
	// сигнал о закрытии контекста будет сигнализировать, что работа сервиса прекращается во время добавления записи (запись добавлена не будет)
	select {
	case <-ctx.Done():
		slog.Error("context want cancel")
		return ctx.Err()
	default:
		_, err := c.OrderRepoCache.GetOrder(order.OrderUID) // проверяем существование заказа (нельзя добавить заказ с одним и тем же id)
		/*
			в switch проверяем: если вернулась ошибка ErrorOrderNotExist, значит записи не существует(и можно добалвять),
			если вернулся nil- запись существует(для нас это ошибка), по умолчанию это любая другая ошибка (запись не добавится)
		*/
		switch {
		case err == storage.ErrorOrderNotExist:
			id, err := c.OrderRepoCache.AddOrder(order)
			if err != nil {
				slog.Error("error add in cache:", "error", err)
				return err
			}
			slog.Info("order add in cache success", "id", id)
			return nil
		case err == nil:
			slog.Error("error add order in cache:", "error", storage.ErrorOrderExist)
			return storage.ErrorOrderExist
		default:
			slog.Error("error add order in cache:", "error", err)
			return err
		}
	}
}

// метод добавления записи в бд
// работа функции аналогична addToCache
func (c Consumer) addToDB(ctx context.Context, order *model.Order) error {
	select {
	case <-ctx.Done():
		slog.Error("context want cancel")
		return ctx.Err()
	default:
		_, err := c.OrderRepoDB.GetOrder(order.OrderUID)
		switch {
		case err == storage.ErrorOrderNotExist:
			id, err := c.OrderRepoDB.AddOrder(order)
			if err != nil {
				slog.Error("error add in db:", "error", err)
				return err
			}
			slog.Info("order add in db success", "id", id)
			return nil
		case err == nil:
			slog.Error("error add order in db:", "error", storage.ErrorOrderExist)
			return storage.ErrorOrderExist
		default:
			slog.Error("error add order in db:", "error", err)
			return err
		}
	}
}
