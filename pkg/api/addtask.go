package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

// Обрабатывает POST-запрос для добавления новой задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("JSON decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid JSON format"})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task title required"})
		return
	}

	if err := checkDate(&task); err != nil {
		log.Printf("date check error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "failed to add task"})
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

// Проверяет и нормализует дату задачи, учитывая правило repeat
func checkDate(task *db.Task) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStr := today.Format(DateFormat)

	// Если дата не указана — ставим сегодняшнюю
	if task.Date == "" {
		task.Date = todayStr
		return nil
	}

	parsedDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("wrong date format: %w", err)
	}

	// Сравнение только по дате, без учёта времени
	if parsedDate.Before(today) {
		if task.Repeat == "" {
			task.Date = todayStr
		} else {
			next, err := NextDate(today, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("repeat rule error: %w", err)
			}
			task.Date = next
		}
	}

	return nil
}

// afterNow возвращает true, если дата now позже даты d (используется для сравнения дат без времени).
func afterNow(now, d time.Time) bool {
	return now.After(d)
}

// writeJSON сериализует переданные данные в JSON и записывает их в ответ клиенту
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
