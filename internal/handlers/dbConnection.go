package handlers

import (
	"database/sql"
	"fmt"
	"sync"
)

type DBConnection struct {
	DB *sql.DB
	Mu sync.Mutex
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
func (c *DBConnection) InitDB() {
	const (
		CreateTableQuery = `CREATE TABLE scheduler (
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		date    CHAR(8) NOT NULL DEFAULT "",
		title   VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT,
		repeat VARCHAR(128) NOT NULL DEFAULT "" 
		);`
	)
	if _, err := c.DB.Exec(CreateTableQuery); err != nil {
		fmt.Println(err)
	}

	if _, err := c.DB.Exec(`CREATE INDEX taks_date ON scheduler (date);`); err != nil {
		fmt.Println(err)

	}

}
