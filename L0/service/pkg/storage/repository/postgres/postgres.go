package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"service/pkg/model"
	"service/pkg/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderPostgresRepository struct {
	Pool *pgxpool.Pool
}

func NewOrderPostgresRepository(pool *pgxpool.Pool) *OrderPostgresRepository {
	return &OrderPostgresRepository{
		Pool: pool,
	}
}

func NewConnPostgres() (*pgxpool.Pool, error) {
	coonString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := pgxpool.New(context.Background(), coonString)
	if err != nil {
		slog.Error("error create connect to db:",
			"error", err,
		)
		return nil, err
	}
	slog.Info("connect to db create success")

	return pool, nil
}

func (p *OrderPostgresRepository) AddOrder(order *model.Order) (string, error) {
	err := InsertOrder(context.Background(), p.Pool, *order)
	if err != nil {
		return "", err
	}
	return order.OrderUID, nil
}

func (p *OrderPostgresRepository) GetOrder(id string) (*model.Order, error) {
	order, err := GetOrderByUID(context.Background(), p.Pool, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Info("Get order: return NoRows by this id", "id", id)
			return nil, storage.ErrorOrderNotExist
		}
		slog.Error("error get order:",
			"id", id,
			"error", err,
		)
		return nil, err
	}
	slog.Info("success get order", "id", id)
	return order, nil
}

func (p *OrderPostgresRepository) DeleteOrder(id string) error {
	err := DeleteOrderByUID(context.Background(), p.Pool, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.ErrorOrderNotExist
		}
		return err
	}
	return nil
}

func (p *OrderPostgresRepository) Close() {
	p.Pool.Close()
}

func (p *OrderPostgresRepository) RestoreCache() (map[string]*model.Order, error) {
	query := `SELECT order_uid FROM orders`

	rows, err := p.Pool.Query(context.Background(), query)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Info("in db no rows, return empty cache")
			return map[string]*model.Order{}, nil
		}
		slog.Error("error restore cache:", "error", err)
		return nil, err
	}

	orders := make(map[string]*model.Order, 0)
	var id string
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			slog.Error("error scan of row:",
				"id", id,
				"error", err)
			return nil, err
		}

		order, err := p.GetOrder(id)
		if err != nil {
			slog.Error("error get order:",
				"id", id,
				"error", err,
			)
			return nil, err
		}

		if _, ok := orders[id]; ok {
			slog.Error("error check order in cache:",
				"id", id,
				"error", storage.ErrorOrderExist,
			)
			return nil, storage.ErrorOrderExist
		}
		orders[id] = order
	}

	return orders, nil
}
