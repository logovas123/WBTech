package kafka

import (
	"context"
	"log/slog"
	"time"

	"service/pkg/storage"

	"github.com/IBM/sarama"
)

/*
Структура Consumer предоставляет собой взаимодействие между связанными компонентами (так как сообщение полученное консьюмером сохраняется в базу и кеш)
*/
type Consumer struct {
	Consumer       sarama.Consumer
	OrderRepoCache storage.OrderRepo
	OrderRepoDB    storage.OrderRepoClose
}

// функция для создания структуры Consumer
func NewConsumer(broker string, cache storage.OrderRepo, db storage.OrderRepoClose) (*Consumer, error) {
	config := sarama.NewConfig()                              // создаём конфиг
	config.Version = sarama.V2_8_2_0                          // задаём версию кафка
	cons, err := sarama.NewConsumer([]string{broker}, config) // создаём консьюмера
	if err != nil {
		slog.Error("error create new consumer:", "error", err)
		return nil, err
	}
	// возвращаем структуру
	return &Consumer{
		Consumer:       cons,
		OrderRepoCache: cache,
		OrderRepoDB:    db,
	}, nil
}

// стартуем консьюмера, чтобы он ждал сообщения
func (c Consumer) StartKafkaConsumer(ctx context.Context) error {
	consumer, err := c.Consumer.ConsumePartition("orders", 0, sarama.OffsetOldest) // создаём консьюмера
	if err != nil {
		slog.Error("error create consumer_partition:",
			"error", err,
		)
		return err
	}
	slog.Info("consumer_partition create success")

	defer consumer.Close() // shutdown consumer

	// в бесконечном цикле ждём сообщения (и обрабатываем его), если пришёл сигнал о закрытии контекста, то выходим из цикла
	for {
		slog.Info("consumer wait new msg...")
		select {
		case msg := <-consumer.Messages():
			slog.Info("get new msg from kafka")
			processCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // контекст с таймаутом - нужен для отсвлеживания слишком долгой обработки сообщения
			defer cancel()

			// обрабатываем сообщение
			if err := c.handleMessage(processCtx, msg.Value); err != nil {
				slog.Error("error handle message:", "error", err)
			}
			slog.Info("msg handle success")

		// сигнал о завершении контекста
		case <-ctx.Done():
			slog.Info("context want cancel")
			return ctx.Err()
		}
	}
}
