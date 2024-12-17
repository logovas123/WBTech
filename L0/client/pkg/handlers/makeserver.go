package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type Server struct {
	mux       *http.ServeMux
	templates *template.Template
}

func NewMuxServer(templates string) *Server {
	srv := new(Server)

	mux := http.NewServeMux()

	mux.HandleFunc("/", srv.IndexPage)

	srv.mux = mux

	srv.templates = template.Must(template.ParseFiles(templates))

	return srv
}

func (s *Server) Start(addr string) error {
	slog.Info("server start")
	if err := http.ListenAndServe(addr, s.mux); err != nil {
		slog.Error("error listen and serve", "error", err)
		return err
	}
	return nil
}
