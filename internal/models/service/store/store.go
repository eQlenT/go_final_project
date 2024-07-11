package store

import (
	"database/sql"
	"errors"
	"fmt"
	"go_final_project/internal/models/service/store/task"
	"time"
)

// Структура TaskStore отвечает за взаимодействие с базой данных для выполнения операций CRUD
// с задачами. Она предоставляет методы для инициализации базы данных, проверки существования задачи по ID,
// удаления задачи, вставки новой задачи, обновления существующей задачи, обновления даты существующей задачи,
// получения всех задач, получения определенной задачи по ID, получения задач, содержащих определенное ключевое слово
// в заголовке или комментарии, и получения задач, запланированных на определенную дату.
type TaskStore struct {
	db *sql.DB
}

// NewTaskStore - это конструктор для структуры TaskStore. Он принимает указатель на экземпляр sql.DB
// в качестве параметра и возвращает новый экземпляр TaskStore.
//
// Параметр db используется для установления соединения с базой данных. Предполагается, что база данных уже создана
// и доступна через предоставленный экземпляр sql.DB.
func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{
		db: db}
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

// CheckID проверяет, существует ли указанный идентификатор в базе данных.
// Он извлекает максимальный идентификатор из таблицы "scheduler" и сравнивает его с указанным идентификатором.
// Если указанный идентификатор больше максимального идентификатора, возвращается ошибка, указывающая, что идентификатор больше числа строк в базе данных.
//
// Параметры:
// - id: Проверяемый идентификатор. Это целое число.
//
// Возвращает:
// - Ошибку, если указанный идентификатор больше максимального идентификатора в базе данных.
// - nil, если указанный идентификатор существует в базе данных.
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

// Delete удаляет задачу из базы данных на основе указанного идентификатора.
//
// Параметры:
// - id: Уникальный идентификатор удаляемой задачи.
//
// Возвращает:
// - Ошибку, если во время удаления произошла ошибка. Если удаление выполнено успешно, возвращается nil
func (s *TaskStore) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

// Insert вставляет новую задачу в базу данных.
//
// Параметры:
// - task: Структура, содержащая данные новой задачи.
//
// Возвращает:
// - Идентификатор вставленной задачи и ошибку, если во время вставки произошла ошибка.
// - Если вставка выполнена успешно, возвращается идентификатор вставленной задачи и nil.
func (s *TaskStore) Insert(task *task.Task) (int, error) {
	res, err := s.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, nil
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// Update обновляет данные существующей задачи в базе данных.
//
// Параметры:
// - task: Структура, содержащая обновленные данные задачи.
//
// Возвращает:
// - Ошибку, если во время обновления произошла ошибка. Если обновление выполнено успешно, возвращается nil.
func (s *TaskStore) Update(task *task.Task) error {
	_, err := s.db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDate обновляет дату существующей задачи в базе данных.
//
// Параметры:
// - task: Структура, содержащая обновленные данные даты задачи.
//
// Возвращает:
// - Ошибку, если во время обновления произошла ошибка. Если обновление выполнено успешно, возвращается nil.
func (s *TaskStore) UpdateDate(task *task.Task) error {
	_, err := s.db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`,
		task.Date, task.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAll извлекает все задачи из базы данных с ограничением на количество возвращаемых записей.
//
// Параметры:
// - limit: Максимальное количество возвращаемых записей.
//
// Возвращает:
// - Словарь, содержащий массив задач и ошибку, если во время извлечения произошла ошибка.
// - Если извлечение выполнено успешно, возвращается словарь с массивом задач и nil.
func (s *TaskStore) GetAll(limit int) (map[string][]task.Task, error) {
	tasks := make(map[string][]task.Task)
	rows, err := s.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
	ORDER BY date LIMIT :limit`,
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := task.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks["tasks"] = append(tasks["tasks"], task)
	}
	if tasks["tasks"] == nil {
		tasks["tasks"] = []task.Task{}
	}
	return tasks, nil
}

// GetTask извлекает конкретную задачу из базы данных на основе указанного идентификатора.
//
// Параметры:
// - id: Уникальный идентификатор извлекаемой задачи.
//
// Возвращает:
// - Указатель на структуру, содержащую данные извлеченной задачи и ошибку, если во время извлечения произошла ошибка.
// - Если извлечение выполнено успешно, возвращается указатель на структуру с данными задачи и nil.
func (s *TaskStore) GetTask(id int) (*task.Task, error) {
	// Получаем задачу по идентификатору
	task := task.Task{}
	rows, err := s.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return &task, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return &task, err
		}
	}
	if err = rows.Err(); err != nil {
		return &task, err
	}
	if task.Date == "" && task.Title == "" && task.Repeat == "" && task.Comment == "" {
		return &task, fmt.Errorf("no rows for id %d", id)
	}
	return &task, nil
}

// GetByWord извлекает все задачи, содержащие указанное ключевое слово в заголовке или комментарии,
// с ограничением на количество возвращаемых записей.
//
// Параметры:
// - key: Ключевое слово для поиска.
// - limit: Максимальное количество возвращаемых записей.
//
// Возвращает:
// - Словарь, содержащий массив задач и ошибку, если во время извлечения произошла ошибка.
// - Если извлечение выполнено успешно, возвращается словарь с массивом задач и nil.
func (s *TaskStore) GetByWord(key string, limit int) (map[string][]task.Task, error) {
	tasks := make(map[string][]task.Task)
	rows, err := s.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
	WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`,
		sql.Named("search", "%"+key+"%"),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := task.Task{}
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
		tasks["tasks"] = []task.Task{}
	}
	return tasks, nil
}

// GetByDate извлекает все задачи, запланированные на указанную дату, с ограничением на количество возвращаемых записей.
//
// Параметры:
// - date: Дата для поиска.
// - limit: Максимальное количество возвращаемых записей.
//
// Возвращает:
// - Словарь, содержащий массив задач и ошибку, если во время извлечения произошла ошибка.
// - Если извлечение выполнено успешно, возвращается словарь с массивом задач и nil.
func (s *TaskStore) GetByDate(date string, limit int) (map[string][]task.Task, error) {
	tasks := make(map[string][]task.Task)
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
		task := task.Task{}
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
		tasks["tasks"] = []task.Task{}
	}
	return tasks, nil
}
