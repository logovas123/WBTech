package main

import (
	"log/slog"
	"os"

	"client/pkg/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env") // читаем env файл
	if err != nil {
		slog.Error("error load .env file:", "error", err)
	}

	pathIndex := "./template/index.html" // путь до темплейта

	srv := handlers.NewMuxServer(pathIndex)
	slog.Info("create new server")

	// стартуем сервер
	err = srv.Start(os.Getenv("CLIENT_HOST") + ":" + os.Getenv("CLIENT_PORT"))
	if err != nil {
		slog.Error("error start server", "error", err)
	}
}
