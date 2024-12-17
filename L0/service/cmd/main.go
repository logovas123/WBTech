package main

import (
	"context"
	"log/slog"

	"service/pkg/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("error parse .env file")
	}
	slog.Info("parse .env file success")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := service.NewService()
	if err != nil {
		slog.Error("error create service:",
			"error", err,
		)
		return
	}
	slog.Info("service create success")

	if err = service.Start(ctx); err != nil {
		slog.Error("error start service:",
			"error", err,
		)
	}

	slog.Info("service gracefull shutdown")
}
