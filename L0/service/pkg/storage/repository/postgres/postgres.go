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

// применяется паттерн репозиторий

// струткура которая содержит пул соединений к базе
/*
При обращение к базе будут выполняться базовые запросы на добавление, получение и удаление заказов(обновление записи не требуется, т.к. сказано, что данные статичны)
*/
type OrderPostgresRepository struct {
	Pool *pgxpool.Pool
}

// создаём новый репозиторий
func NewOrderPostgresRepository(pool *pgxpool.Pool) *OrderPostgresRepository {
	return &OrderPostgresRepository{
		Pool: pool,
	}
}

// функция для создания пула соединений до базы
func NewConnPostgres() (*pgxpool.Pool, error) {
	// создаём адрес подключения к базе
	coonString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	// создаём пул соединений
	pool, err := pgxpool.New(context.Background(), coonString)
	if err != nil {
		slog.Error("error create connect to db:",
			"error", err,
		)
		return nil, err
	}
	slog.Info("connect to db create success")

	return pool, nil // возвращаем пул соединений
}

// метод для добавления заказа
func (p *OrderPostgresRepository) AddOrder(order *model.Order) (string, error) {
	err := InsertOrder(context.Background(), p.Pool, *order)
	if err != nil {
		return "", err
	}
	return order.OrderUID, nil
}

// метод для получения заказа
func (p *OrderPostgresRepository) GetOrder(id string) (*model.Order, error) {
	order, err := GetOrderByUID(context.Background(), p.Pool, id) // получаем заказ
	if err != nil {
		// если вернулся pgx.ErrNoRows, возвращаем кастомную ошибку об отсутствии записи
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
	return order, nil // возвращаем заказ
}

// метод для удаления заказа
func (p *OrderPostgresRepository) DeleteOrder(id string) error {
	err := DeleteOrderByUID(context.Background(), p.Pool, id) // удаляем заказ
	if err != nil {
		// при отсутствии записи возвращаем ошибку
		if err == sql.ErrNoRows {
			return storage.ErrorOrderNotExist
		}
		return err
	}
	return nil
}

// метод для закрытия пула соединений
func (p *OrderPostgresRepository) Close() {
	p.Pool.Close()
}

/*
Метод для востановления кеша.
Сначала выполняется запрос на получение всех id в бд. Затем по каждому id мы получаем заказ из бд и сохраняем в кеш.
*/
func (p *OrderPostgresRepository) RestoreCache() (map[string]*model.Order, error) {
	query := `SELECT order_uid FROM orders`

	rows, err := p.Pool.Query(context.Background(), query) // получаем записи которые хранят id
	if err != nil {
		// если пришла данная ошибка, то это говорит о том, что бд была пустой
		if err == pgx.ErrNoRows {
			slog.Info("in db no rows, return empty cache")
			return map[string]*model.Order{}, nil
		}
		slog.Error("error restore cache:", "error", err)
		return nil, err
	}

	orders := make(map[string]*model.Order, 0)
	var id string
	// получаем записи по id
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

	return orders, nil // возвращаем восстановленный кеш
}
