package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"service/pkg/handler"
	"service/pkg/kafka"
	"service/pkg/storage/repository/cache"
	"service/pkg/storage/repository/postgres"
)

// в пакете service определяем структуру сервиса и её методы
// структура сервиса будет состоять из хендлера, который будет обрабатыввать внешние запросы
// и консьюмера, который будет получать сообщения от брокера
type Service struct {
	Handler  http.Handler
	Consumer *kafka.Consumer
}

// функция для создания нового сервиса
func NewService() (*Service, error) {
	cache := cache.NewCashe() // создаём кеш
	slog.Info("cache create success")

	pool, err := postgres.NewConnPostgres() // создаём новый коннект к бд
	if err != nil {
		slog.Error("error create coonect to db:",
			"error", err,
		)
		return nil, err
	}
	slog.Info("new connect to db create success")

	db := postgres.NewOrderPostgresRepository(pool) // создаём базу
	slog.Info("new db repository create success")

	slog.Info("attempt to restore cache from db...")
	cache.Orders, err = db.RestoreCache() // восстанавливаем кеш
	if err != nil {
		slog.Error("error restore cache:",
			"error", err,
		)
		return nil, err
	}
	slog.Info("cache restore success")

	consumer, err := kafka.NewConsumer(os.Getenv("BROKER_HOST")+":"+os.Getenv("BROKER_PORT"), cache, db) // создаём нового консьюмера, который будет принимать сообщения от брокера
	if err != nil {
		slog.Error("error create new consumer:", "error", err)
		return nil, err
	}
	slog.Info("create new consumer success")

	// создаём хендлер для обработки запросов
	OrderHandler := &handler.OrderHandler{
		OrderRepo: cache,
	}
	slog.Info("create new handler success")

	mux := handler.NewMuxServer(OrderHandler) // создаём роутер
	slog.Info("create new router success")

	// возвращаем структуру сервиса
	return &Service{
		Handler:  mux,
		Consumer: consumer,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM) // создаём копию контекста, который будет закрыт, когда  придёт один из перечисленных сигналов
	defer stop()

	// запускаем параллельный процесс:
	// консьюмер ждёт когда придёт сообщение
	go func() {
		err := s.Consumer.StartKafkaConsumer(ctx) // стартуем консьюмера
		if err != nil {
			slog.Error("error in kafka:", "error", err)
		}
		slog.Info("kafka stopped success")
	}()

	srv := http.Server{}                                                  // создаём сервер
	srv.Addr = os.Getenv("HOST_SERVICE") + ":" + os.Getenv("SERVER_PORT") // адрес сервера
	if srv.Addr == "" {
		srv.Addr = "localhost:8088"
	}
	srv.Handler = s.Handler // передали хендлер

	// запускаем параллельный процесс:
	// сервер слушает соединение
	go func() {
		slog.Info("server start")
		// если придёт ErrServerClosed, то это говорит о корректном завершении рбаоты сервиса
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error listen server:", "error", err)
		}
	}()

	<-ctx.Done() // здесь функция блокируется, пока не придёт сигнал о завершении

	// начинаем корректно завершать работу сервиса

	s.Consumer.OrderRepoDB.Close() // закрываем пул соединений к бд

	slog.Info("db pool success closed")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second) // создаём контекст с таймаутом, если не удасться корректно завершить работу сервера(слишком долго)
	defer cancel()
	// корректно завершаем работу сервиса
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("error shutdown server:", "error", err)
		return err
	}

	slog.Info("server stopped success")

	return nil
}
