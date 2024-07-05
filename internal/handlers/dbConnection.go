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
