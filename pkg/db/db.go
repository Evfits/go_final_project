package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB // Глобальная переменная для использования БД в других пакетах

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// Init инициализирует базу данных и создаёт таблицу, если файл ещё не существует
func Init(dbFile string) error {
	// Проверка наличия файла
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открытие БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть БД: %w", err)
	}

	// Проверка соединения
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ошибка соединения с БД: %w", err)
	}

	// Если файл создавался заново — выполняем schema
	if install {
		_, err := DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания схемы БД: %w", err)
		}
	}

	return nil
}

// Close Закрывает соединение с БД (если открыто), возвращает ошибку при неудаче
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
