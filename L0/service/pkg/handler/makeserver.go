package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewMuxServer(orderHandler *OrderHandler) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/{id}", orderHandler.GetOrderByID).Methods(http.MethodGet)

	return r
}
