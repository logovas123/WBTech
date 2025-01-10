package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

/*
При методе get выводится главная страница(без данных)
При методе post (запрос выполняется по нажатию на кнопку) обработчик получает данные из формы(id).
Делает запрос к сервису, получает в ответ либо ошибку, либо данные заказа. Применяет полученные данные к шаблону
*/

func (s *Server) IndexPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := s.templates.Execute(w, nil) // применяем шаблон
		if err != nil {
			slog.Error("error execute template", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		id := r.FormValue("id")                                                                                                              // получаем данные формы
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/api/%s", os.Getenv("HOST_SERVICE"), os.Getenv("SERVER_PORT"), id), nil) // создаём запрос к сервису
		if err != nil {
			slog.Error("error create request", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req) // выполняем запрос
		if err != nil {
			slog.Error("error complete request", "error", err)
			err = s.templates.Execute(w, "error request in external service")
			if err != nil {
				slog.Error("error execute template", "error", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			return
		}
		defer resp.Body.Close() // закрываем тело ответа

		bodyResp, err := io.ReadAll(resp.Body) // читаем тело ответа
		if err != nil {
			slog.Error("error read response body")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)                                         // устанавливаем статус
		err = s.templates.ExecuteTemplate(w, "index.html", string(bodyResp)) // применяем данные к шаблону
		if err != nil {
			slog.Error("error execute template", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

	// другой метод не зарегистрирован
	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}
