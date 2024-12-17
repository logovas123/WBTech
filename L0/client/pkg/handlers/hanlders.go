package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func (s *Server) IndexPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := s.templates.Execute(w, nil)
		if err != nil {
			slog.Error("error execute template", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		id := r.FormValue("id")
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/api/%s", os.Getenv("HOST_SERVICE"), os.Getenv("SERVER_PORT"), id), nil)
		if err != nil {
			slog.Error("error create request", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req)
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
		defer resp.Body.Close()

		bodyResp, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("error read response body")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = s.templates.ExecuteTemplate(w, "index.html", string(bodyResp))
		if err != nil {
			slog.Error("error execute template", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}
