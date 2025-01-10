package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

// структура для гибкого управления сервером
type Server struct {
	mux       *http.ServeMux
	templates *template.Template
}

// функция создаёт новую структур Server
func NewMuxServer(templates string) *Server {
	srv := new(Server)

	mux := http.NewServeMux() // создаём роутер

	mux.HandleFunc("/", srv.IndexPage) // регистрируем обработчик

	srv.mux = mux

	srv.templates = template.Must(template.ParseFiles(templates)) // парсим файл html, возвраащем темплейт

	return srv // возвращаем сервер
}

// метод для старта сервера
func (s *Server) Start(addr string) error {
	slog.Info("server start")
	if err := http.ListenAndServe(addr, s.mux); err != nil {
		slog.Error("error listen and serve", "error", err)
		return err
	}
	return nil
}
