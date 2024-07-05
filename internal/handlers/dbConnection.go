package handlers

import (
	"database/sql"
	"sync"
)

type DBConnection struct {
	DB *sql.DB
	Mu sync.Mutex
}
