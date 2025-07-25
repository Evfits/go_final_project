package server

import (
	"fmt"
	"go_final_project/pkg/api"
	"log"
	"net/http"
	"os"
)

// Start запускает HTTP-сервер с обработкой API и статических файлов.
// webDir — путь к фронтенду (директория с index.html и js/css).
func Start(webDir string) {
	// Получение порта из переменной окружения/значение по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Создание маршрутизатора
	mux := http.NewServeMux()

	// Инициализация API-роутов
	api.RegisterRoutes(mux)

	// Добавление файлового сервера для фронтенда (index.html, js, css, и т.д.)
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Сервер запущен: http://localhost%s/ (статичные файлы из %s)", addr, webDir)

	// Запуск сервера
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
