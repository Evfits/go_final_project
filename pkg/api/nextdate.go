package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NextDate рассчитывает следующую дату по правилу repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неправильный формат даты dstart: %w", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat не указан")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("некорректное правило repeat")
	}

	switch parts[0] {
	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("отсутствует количество дней")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", fmt.Errorf("некорректный интервал: %s", parts[1])
		}
		for {
			startDate = startDate.AddDate(0, 0, interval)
			if startDate.After(now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	case "y":
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if startDate.After(now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемое правило: %s", parts[0])
	}
}

// nextDateHandler обрабатывает запросы и возвращает следующую дату задачи
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	nowStr := query.Get("now")
	dstart := query.Get("date")
	repeat := query.Get("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат now", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
