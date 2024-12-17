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
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("error load .env file:", "error", err)
	}
	broker := os.Getenv("BROKER_HOST") + ":" + os.Getenv("BROKER_PORT")
	producer, err := kafka.NewProducer(broker)
	if err != nil {
		slog.Error("error of create producer", "error", err)
		return
	}
	defer producer.Close()

	slog.Info("start new producer")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	topic := "orders"
	key := "order_key"

	err = producer.SendMessage(ctx, topic, key)
	if err != nil {
		slog.Error("error to send messages", "error", err)
	}

	slog.Info("producer closed")
}
