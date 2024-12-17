package kafka

import (
	"context"
	"log/slog"
	"time"

	"service/pkg/storage"

	"github.com/IBM/sarama"
)

type Consumer struct {
	Consumer       sarama.Consumer
	OrderRepoCache storage.OrderRepo
	OrderRepoDB    storage.OrderRepoClose
}

func NewConsumer(broker string, cache storage.OrderRepo, db storage.OrderRepoClose) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_2_0
	cons, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		slog.Error("error create new consumer:", "error", err)
		return nil, err
	}
	return &Consumer{
		Consumer:       cons,
		OrderRepoCache: cache,
		OrderRepoDB:    db,
	}, nil
}

func (c Consumer) StartKafkaConsumer(ctx context.Context) error {
	consumer, err := c.Consumer.ConsumePartition("orders", 0, sarama.OffsetOldest)
	if err != nil {
		slog.Error("error create consumer_partition:",
			"error", err,
		)
		return err
	}
	slog.Info("consumer_partition create success")

	defer consumer.Close()
	for {
		slog.Info("consumer wait new msg...")
		select {
		case msg := <-consumer.Messages():
			slog.Info("get new msg from kafka")
			processCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := c.handleMessage(processCtx, msg.Value); err != nil {
				slog.Error("error handle message:", "error", err)
			}
			slog.Info("msg handle success")

		case <-ctx.Done():
			slog.Info("context want cancel")
			return ctx.Err()
		}
	}
}
