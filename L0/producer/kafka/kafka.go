package kafka

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"producer/model"

	"github.com/IBM/sarama"
	"github.com/brianvoe/gofakeit/v7"
)

type Producer struct {
	Producer sarama.SyncProducer
}

func NewProducer(broker string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 2
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_8_2_0
	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return nil, err
	}
	return &Producer{Producer: producer}, nil
}

func (p *Producer) SendMessage(ctx context.Context, topic string, key string) error {
	scanner := bufio.NewScanner(os.Stdin)
	var order model.Order
	slog.Info("\nPlease, enter \"send\" for sending msg.")
	for scanner.Scan() {
		if scanner.Text() != "send" {
			slog.Info("This is not \"send\"")
			fmt.Println()
			slog.Info("Please, enter \"send\" for sending msg.")
			continue
		}

		gofakeit.Struct(&order)
		value, err := json.Marshal(order)
		if err != nil {
			slog.Error("error marshal order", "error", err)
			return err
		}

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.ByteEncoder(value),
		}

		_, _, err = p.Producer.SendMessage(msg)
		if err != nil {
			slog.Error("error send message", "error", err)
			return err
		}

		slog.Info("message sent successfully", "topic", topic, "id", order.OrderUID)
		fmt.Println()
		slog.Info("Please, enter \"send\" for sending msg.")
	}

	return nil
}

func (p *Producer) Close() error {
	return p.Producer.Close()
}
