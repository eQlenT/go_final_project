package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB, logger *zap.SugaredLogger) *TaskStore {
	return &TaskStore{db: db}
}

// InitDB инициализирует базу данных путем создания таблицы "scheduler" и индекса на столбце "date".
// Таблица "scheduler" имеет следующие столбцы:
// - id: первичный ключ целого числа с автоинкрементом
// - date: поле символьного типа длиной 8 символов, не может быть NULL, имеет значение по умолчанию пустая строка
// - title: поле переменной длины строки с максимальной длиной 128 символов, не может быть NULL, имеет значение по умолчанию пустая строка
// - comment: текстовое поле для хранения комментариев
// - repeat: поле переменной длины строки с максимальной длиной 128 символов, не может быть NULL, имеет значение по умолчанию пустая строка
//
// Если во время выполнения SQL-запросов возникает ошибка, она будет выведена в консоль.
func (s *TaskStore) InitDB() {
	const (
		CreateTableQuery = `CREATE TABLE scheduler (
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		date    CHAR(8) NOT NULL DEFAULT "",
		title   VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT,
		repeat VARCHAR(128) NOT NULL DEFAULT "" 
		);`
	)
	if _, err := s.db.Exec(CreateTableQuery); err != nil {
		fmt.Println(err)
	}

	if _, err := s.db.Exec(`CREATE INDEX taks_date ON scheduler (date);`); err != nil {
		fmt.Println(err)

	}

}

func (s *TaskStore) CheckID(id int) error {
	var maxID int
	row := s.db.QueryRow(`SELECT MAX(id) FROM scheduler`)
	row.Scan(&maxID)
	if err := row.Err(); err != nil {
		return err
	}
	if id > maxID {
		err := errors.New("given id is more than number of rows in DB")
		return err
	}
	return nil
}

func (s *TaskStore) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskStore) Insert(task *Task) (int, error) {
	res, err := s.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, nil
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (s *TaskStore) Update(task *Task) error {
	_, err := s.db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskStore) UpdateDate(task *Task) error {
	_, err := s.db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`,
		task.Date, task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskStore) GetAll(limit int) (map[string][]Task, error) {
	tasks := make(map[string][]Task)
	rows, err := s.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
	ORDER BY date LIMIT :limit`,
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks["tasks"] = append(tasks["tasks"], task)
	}
	if tasks["tasks"] == nil {
		tasks["tasks"] = []Task{}
	}
	return tasks, nil
}

func (s *TaskStore) GetTask(id int) (*Task, error) {
	// Получаем задачу по идентификатору
	task := Task{}
	rows, err := s.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return &Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return &Task{}, err
		}
	}
	if err = rows.Err(); err != nil {
		return &Task{}, err
	}
	if task.Date == "" && task.Title == "" && task.Repeat == "" && task.Comment == "" {
		return &Task{}, fmt.Errorf("no rows for id %d", id)
	}
	return &task, nil
}

func (s *TaskStore) GetByWord(key string, limit int) (map[string][]Task, error) {
	tasks := make(map[string][]Task)
	rows, err := s.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
	WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`,
		sql.Named("search", "%"+key+"%"),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks["tasks"] = append(tasks["tasks"], task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if tasks["tasks"] == nil {
		tasks["tasks"] = []Task{}
	}
	return tasks, nil
}

func (s *TaskStore) GetByDate(date string, limit int) (map[string][]Task, error) {
	tasks := make(map[string][]Task)
	dateTime, err := time.Parse("02.01.2006", date)
	if err != nil {
		return nil, err
	}
	dateFormat := dateTime.Format("20060102")
	rows, err := s.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
		WHERE date = :date LIMIT :limit`,
		sql.Named("date", dateFormat),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks["tasks"] = append(tasks["tasks"], task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if tasks["tasks"] == nil {
		tasks["tasks"] = []Task{}
	}
	return tasks, nil
}
