package handler

import (
	"encoding/json"
	"net/http"

	"service/pkg/storage"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	OrderRepo storage.OrderRepo
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	order, err := h.OrderRepo.GetOrder(id)
	if err != nil {
		if err == storage.ErrorOrderNotExist {
			http.Error(w, "order with this id not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}

	body, err := json.MarshalIndent(*order, "", "    ")
	if err != nil {
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(body)
}
