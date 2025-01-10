package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"producer/kafka"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env") // читаем env файл
	if err != nil {
		slog.Error("error load .env file:", "error", err)
	}
	broker := os.Getenv("BROKER_HOST") + ":" + os.Getenv("BROKER_PORT") // получаем из переменных окружения значения хоста и порта, и создаём адрес брокера
	producer, err := kafka.NewProducer(broker)                          // создаём продюсера
	if err != nil {
		slog.Error("error of create producer", "error", err)
		return
	}
	defer producer.Close() // shutdown producer

	slog.Info("start new producer")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // создаём контекст с таймаутом (нужен если операция будет долго выполняться)
	defer cancel()

	topic := "orders"
	key := "order_key"

	err = producer.SendMessage(ctx, topic, key) // отправляем сообщение
	if err != nil {
		slog.Error("error to send messages", "error", err)
	}

	slog.Info("producer closed")
}
