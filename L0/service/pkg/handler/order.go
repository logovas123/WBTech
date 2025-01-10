package handler

import (
	"encoding/json"
	"net/http"

	"service/pkg/storage"

	"github.com/gorilla/mux"
)

// структкура хендлера, содержит только кеш, так как взаимодействует только с кешем
type OrderHandler struct {
	OrderRepo storage.OrderRepo
}

// обработчик запроса который получает из url запроса id, по id получает заказ из кеша, маршалит его и пишет в тело ответа
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // получаем параметры запроса (свойственно только gorilla/mux)

	id := vars["id"] // получаем id

	order, err := h.OrderRepo.GetOrder(id) // получаем заказ из кеша, если вернулась ErrorOrderNotExist, возвращаем запись о том, что заказ не найден
	if err != nil {
		if err == storage.ErrorOrderNotExist {
			http.Error(w, "order with this id not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}

	body, err := json.MarshalIndent(*order, "", "    ") // маршалим полученный заказ(MarshalIndent нужен для красивого отбражения во фронте)
	if err != nil {
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}

	// установили заголовки, установили статус, записали тело ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(body)
}
