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

// пакет kafka нужен для определения структуры Producer и описания его методов
// создаём структуру Producer
type Producer struct {
	Producer sarama.SyncProducer
}

// метод создаёт новую струтктуру Producer
func NewProducer(broker string) (*Producer, error) {
	config := sarama.NewConfig() // создаём конфиг для кафка
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 2 // число ретраев
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_8_2_0                                  // версия кафки
	producer, err := sarama.NewSyncProducer([]string{broker}, config) // создаём продюсера
	if err != nil {
		return nil, err
	}
	return &Producer{Producer: producer}, nil
}

/*
Метод нужен для отправки сообщений.
Действие метода такое: при вводе в консоль слова send, генерируются случайные данные на основе модели, которая лежит в папке model/
Формируется сообщение и отправляется в канал
*/
func (p *Producer) SendMessage(ctx context.Context, topic string, key string) error {
	scanner := bufio.NewScanner(os.Stdin) // создаём сканер, в бесконечном цикле читаем значения, ожидая слова send
	var order model.Order                 // струткура для отправки
	slog.Info("\nPlease, enter \"send\" for sending msg.")
	for scanner.Scan() {
		if scanner.Text() != "send" {
			slog.Info("This is not \"send\"")
			fmt.Println()
			slog.Info("Please, enter \"send\" for sending msg.")
			continue
		}

		gofakeit.Struct(&order)           // генерируем случайную структуру
		value, err := json.Marshal(order) // маршалим
		if err != nil {
			slog.Error("error marshal order", "error", err)
			return err
		}

		// создаём сообщение, передавая в неё топик, ключ и значение (преобразую к нужному типу)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.ByteEncoder(value),
		}

		_, _, err = p.Producer.SendMessage(msg) // отправляем сообщение
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
	return p.Producer.Close() // shutdown producer
}
