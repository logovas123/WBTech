package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewMuxServer(orderHandler *OrderHandler) http.Handler {
	r := mux.NewRouter() // создаём роутер

	r.HandleFunc("/api/{id}", orderHandler.GetOrderByID).Methods(http.MethodGet) // регистрируем обработчик

	return r
}
