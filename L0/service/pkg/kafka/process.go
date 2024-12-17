package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"service/pkg/model"
	"service/pkg/storage"
)

func (c Consumer) handleMessage(ctx context.Context, message []byte) error {
	order := &model.Order{}

	err := json.Unmarshal(message, order)
	if err != nil {
		slog.Error("error unmarshal message:", "error", err)
		return err
	}
	slog.Info("msg unmarshal success", "order id", order.OrderUID)

	err = c.addToCache(ctx, order)
	if err != nil {
		slog.Error("error add to cache")
		return err
	}
	slog.Info("order add in cache success")

	err = c.addToDB(ctx, order)
	if err != nil {
		slog.Error("error add to db", "error", err)
		return err
	}
	slog.Info("order add in db success")

	return nil
}

func (c Consumer) addToCache(ctx context.Context, order *model.Order) error {
	select {
	case <-ctx.Done():
		slog.Error("context want cancel")
		return ctx.Err()
	default:
		_, err := c.OrderRepoCache.GetOrder(order.OrderUID)
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
