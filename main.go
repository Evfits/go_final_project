package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	// Получение пути к базе данных из переменной окружения или по умолчанию
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализация базы данных
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Не удалось инициализировать БД: %v", err)
	}
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}
	}()

	// Запуск сервера, указывая директорию с фронтендом
	server.Start("./web")
}
