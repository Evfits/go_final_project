package api

import "net/http"

// RegisterRoutes регистрирует обработчики API
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.HandleFunc("/api/task/done", doneTaskHandler)
}
