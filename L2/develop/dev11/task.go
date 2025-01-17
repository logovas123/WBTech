package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"res": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

// структура события
type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
}

// Calendar - хранилище событий
type Calendar struct {
	events map[int]Event
	nextID int
}

func NewCalendar() *Calendar {
	return &Calendar{
		events: make(map[int]Event),
		nextID: 1,
	}
}

// создаём эвент
func (c *Calendar) CreateEvent(userID int, title string, date time.Time) Event {
	event := Event{
		ID:     c.nextID,
		UserID: userID,
		Title:  title,
		Date:   date,
	}
	c.events[c.nextID] = event
	c.nextID++
	return event
}

// обновление события
func (c *Calendar) UpdateEvent(eventID int, userID int, title string, date time.Time) error {
	if _, exists := c.events[eventID]; !exists {
		return errors.New("event not found")
	}
	c.events[eventID] = Event{
		ID:     eventID,
		UserID: userID,
		Title:  title,
		Date:   date,
	}
	return nil
}

// удаление события
func (c *Calendar) DeleteEvent(eventID int) error {
	if _, exists := c.events[eventID]; !exists {
		return errors.New("event not found")
	}
	delete(c.events, eventID)
	return nil
}

func (c *Calendar) GetEvents(start, end time.Time) []Event {
	var result []Event
	for _, event := range c.events {
		if event.Date.After(start) && event.Date.Before(end) {
			result = append(result, event)
		}
	}
	return result
}

// Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func parseFormInt(values map[string][]string, key string) (int, error) {
	if val, ok := values[key]; ok && len(val) > 0 {
		return strconv.Atoi(val[0])
	}
	return 0, errors.New("missing or invalid " + key)
}

func parseFormDate(values map[string][]string, key string) (time.Time, error) {
	if val, ok := values[key]; ok && len(val) > 0 {
		return time.Parse("2006-01-02", val[0])
	}
	return time.Time{}, errors.New("missing or invalid " + key)
}

func handleCreateEvent(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
			return
		}

		userID, err := parseFormInt(r.Form, "user_id")
		if err != nil {
			http.Error(w, `{"error": "invalid user_id"}`, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		if title == "" {
			http.Error(w, `{"error": "title is required"}`, http.StatusBadRequest)
			return
		}

		date, err := parseFormDate(r.Form, "date")
		if err != nil {
			http.Error(w, `{"error": "invalid date"}`, http.StatusBadRequest)
			return
		}

		event := calendar.CreateEvent(userID, title, date)
		response, _ := json.Marshal(map[string]interface{}{"result": event})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func handleUpdateEvent(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
			return
		}

		eventID, err := parseFormInt(r.Form, "event_id")
		if err != nil {
			http.Error(w, `{"error": "invalid event_id"}`, http.StatusBadRequest)
			return
		}

		userID, err := parseFormInt(r.Form, "user_id")
		if err != nil {
			http.Error(w, `{"error": "invalid user_id"}`, http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		if title == "" {
			http.Error(w, `{"error": "title is required"}`, http.StatusBadRequest)
			return
		}

		date, err := parseFormDate(r.Form, "date")
		if err != nil {
			http.Error(w, `{"error": "invalid date"}`, http.StatusBadRequest)
			return
		}

		if err := calendar.UpdateEvent(eventID, userID, title, date); err != nil {
			http.Error(w, `{"error": "event not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "updated"}`))
	}
}

func handleDeleteEvent(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
			return
		}

		eventID, err := parseFormInt(r.Form, "event_id")
		if err != nil {
			http.Error(w, `{"error": "invalid event_id"}`, http.StatusBadRequest)
			return
		}

		if err := calendar.DeleteEvent(eventID); err != nil {
			http.Error(w, `{"error": "event not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "deleted"}`))
	}
}

func handleEventsForPeriod(calendar *Calendar, days int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		date, err := parseFormDate(query, "date")
		if err != nil {
			http.Error(w, `{"error": "invalid date"}`, http.StatusBadRequest)
			return
		}

		start := date
		end := date.AddDate(0, 0, days)
		events := calendar.GetEvents(start, end)

		response, _ := json.Marshal(map[string]interface{}{"result": events})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func main() {
	calendar := NewCalendar()

	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", handleCreateEvent(calendar))
	mux.HandleFunc("/update_event", handleUpdateEvent(calendar))
	mux.HandleFunc("/delete_event", handleDeleteEvent(calendar))
	mux.HandleFunc("/events_for_day", handleEventsForPeriod(calendar, 1))
	mux.HandleFunc("/events_for_week", handleEventsForPeriod(calendar, 7))
	mux.HandleFunc("/events_for_month", handleEventsForPeriod(calendar, 30))

	loggedMux := loggingMiddleware(mux)

	port := "8080"
	fmt.Println("Server start")
	log.Fatal(http.ListenAndServe(":"+port, loggedMux))
}
