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

type Service struct {
	Handler  http.Handler
	Consumer *kafka.Consumer
}

func NewService() (*Service, error) {
	cache := cache.NewCashe()
	slog.Info("cache create success")

	pool, err := postgres.NewConnPostgres()
	if err != nil {
		slog.Error("error create coonect to db:",
			"error", err,
		)
		return nil, err
	}
	slog.Info("new connect to db create success")

	db := postgres.NewOrderPostgresRepository(pool)
	slog.Info("new db repository create success")

	slog.Info("attempt to restore cache from db...")
	cache.Orders, err = db.RestoreCache()
	if err != nil {
		slog.Error("error restore cache:",
			"error", err,
		)
		return nil, err
	}
	slog.Info("cache restore success")

	consumer, err := kafka.NewConsumer(os.Getenv("BROKER_HOST")+":"+os.Getenv("BROKER_PORT"), cache, db)
	if err != nil {
		slog.Error("error create new consumer:", "error", err)
		return nil, err
	}
	slog.Info("create new consumer success")

	OrderHandler := &handler.OrderHandler{
		OrderRepo: cache,
	}
	slog.Info("create new handler success")

	mux := handler.NewMuxServer(OrderHandler)
	slog.Info("create new router success")

	return &Service{
		Handler:  mux,
		Consumer: consumer,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err := s.Consumer.StartKafkaConsumer(ctx)
		if err != nil {
			slog.Error("error in kafka:", "error", err)
		}
		slog.Info("kafka stopped success")
	}()

	srv := http.Server{}
	srv.Addr = os.Getenv("HOST_SERVICE") + ":" + os.Getenv("SERVER_PORT")
	if srv.Addr == "" {
		srv.Addr = "localhost:8088"
	}
	srv.Handler = s.Handler

	go func() {
		slog.Info("server start")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error listen server:", "error", err)
		}
	}()

	<-ctx.Done()

	s.Consumer.OrderRepoDB.Close()

	slog.Info("db pool success closed")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("error shutdown server:", "error", err)
		return err
	}

	slog.Info("server stopped success")

	return nil
}
