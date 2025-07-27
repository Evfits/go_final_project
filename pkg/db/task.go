package db

import (
	"fmt"
	"strings"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в таблицу, возвращает идентификатор созданной записи
func AddTask(task *Task) (int64, error) {
	query := `
INSERT INTO scheduler (date, title, comment, repeat)
VALUES (?, ?, ?, ?)
`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Tasks возвращает список задач, отсортированный по дате, с возможностью поиска
func Tasks(limit int, search string) ([]*Task, error) {
	baseQuery := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
    `
	var args []interface{}
	where := ""
	if len(search) == len("02.01.2006") && strings.Count(search, ".") == 2 {
		// Поиск по дате: 02.01.2006 → 20060102
		parts := strings.Split(search, ".")
		if len(parts) == 3 {
			search = parts[2] + parts[1] + parts[0]
			where = " WHERE date = ?"
			args = append(args, search)
		}
	} else if search != "" {
		// Поиск по заголовку и комментарию
		where = " WHERE title LIKE ? OR comment LIKE ?"
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	query := baseQuery + where + " ORDER BY date ASC LIMIT ?"
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

// GetTask возвращает задачу по идентификатору id
func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?`, id)
	var t Task
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTask обновляет все поля задачи по id
func UpdateTask(t *Task) error {
	res, err := DB.Exec(`UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`,
		t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Задача не найдена")
	}
	return nil
}

// DeleteTask удаляет задачу по идентификатору
func DeleteTask(id string) error {
	var idInt int64
	if _, err := fmt.Sscan(id, &idInt); err != nil {
		return fmt.Errorf("неверный формат id")
	}

	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, idInt)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удалённых строк: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("задача с таким id не найдена")
	}
	return nil
}

func UpdateDate(date string, id string) error {
	_, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date, id)
	return err
}
