package postgres

import (
	"context"
	"log/slog"

	"service/pkg/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertOrder(ctx context.Context, pool *pgxpool.Pool, order model.Order) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	orderQuery := `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.Exec(ctx, orderQuery,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	deliveryQuery := `
		INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(ctx, deliveryQuery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.ZIP,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return err
	}

	paymentQuery := `
		INSERT INTO payment (
			order_uid, transaction_id, request_id, currency, provider, amount, payment_dt,
			bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.Exec(ctx, paymentQuery,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	itemQuery := `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price,
			nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemQuery,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func GetOrderByUID(ctx context.Context, pool *pgxpool.Pool, orderUID string) (*model.Order, error) {
	var order model.Order

	err := pool.QueryRow(ctx, `
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		slog.Error("error get order, error in Scan():",
			"error", err,
		)
		return nil, err
	}

	err = pool.QueryRow(ctx, `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1
	`, orderUID).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.ZIP,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		slog.Error("error get order, error in Scan():",
			"error", err,
		)
		return nil, err
	}

	err = pool.QueryRow(ctx, `
		SELECT transaction_id, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1
	`, orderUID).Scan(
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDT,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		slog.Error("error get order, error in Scan():",
			"error", err,
		)
		return nil, err
	}

	rows, err := pool.Query(ctx, `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`, orderUID)
	if err != nil {
		slog.Error("error get order, error in Scan():",
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err = rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			slog.Error("error get order, error in Scan():",
				"error", err,
			)
			return nil, err
		}

		order.Items = append(order.Items, item)
	}

	if rows.Err() != nil {
		slog.Error("error get order, error in rows.Err():",
			"id", orderUID,
			"error", rows.Err(),
		)
		return nil, rows.Err()
	}

	slog.Info("order get success", "id", orderUID)

	return &order, nil
}

func DeleteOrderByUID(ctx context.Context, pool *pgxpool.Pool, orderUID string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM orders WHERE order_uid = $1`
	_, err = pool.Exec(ctx, query, orderUID)
	if err != nil {
		return err
	}

	return nil
}
