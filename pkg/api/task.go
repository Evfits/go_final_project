package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запрос /api/tasks и возвращает список ближайших задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		log.Printf("failed to get tasks: %v", err)
		http.Error(w, `{"error":"failed to get tasks"}`, http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = []*db.Task{} // избегаем null в JSON
	}
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

// getTaskHandler обрабатывает GET-запрос и возвращает параметры задачи.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task id is required"})
		return
	}
	t, err := db.GetTask(id)
	if err != nil {
		log.Printf("task not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]string{"error": "task not found"})
		return
	}
	writeJSON(w, t)
}

// editTaskHandler обрабатывает PUT-запрос /api/task и редактирует задачу по id
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Printf("json decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid json format"})
		return
	}
	if t.ID == "" || t.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid task data"})
		return
	}

	if err := checkDate(&t); err != nil {
		log.Printf("date check error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&t); err != nil {
		log.Printf("update task error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{})
}

// doneTaskHandler обрабатывает POST-запрос /api/task/done и отмечает задачу выполненной с учётом repeat
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task id is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("task not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]string{"error": "task not found"})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			log.Printf("delete task error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		log.Printf("repeat rule error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		log.Printf("update date error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}

// deleteTaskHandler обрабатывает DELETE-запрос /api/task и удаляет задачу по id
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task id is required"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		log.Printf("delete task error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{}) // Пустой JSON — всё прошло успешно
}

// taskHandler обрабатывает запросы /api/task для всех поддерживаемых методов: POST, GET, PUT, DELETE.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
