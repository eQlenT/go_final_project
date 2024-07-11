package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// CheckDB проверяет, существует ли файл базы данных SQLite и возвращает его путь.
// Если файла не существует, возвращается имя файла базы данных по умолчанию и флаг, указывающий, что базу данных необходимо установить.
//
// Возвращает:
// path (string): Путь к файлу базы данных SQLite.
// install (bool): Флаг, указывающий, что базу данных необходимо установить.
func CheckDB() (string, bool) {

	// Получение пути к файлу базы данных SQLite из переменной окружения TODO_DBFILE.
	// Если переменная окружения не установлена, используется имя файла базы данных по умолчанию "scheduler.db".
	path := os.Getenv("TODO_DBFILE")
	if path == "" {
		path = "scheduler.db"
	}

	// Получение абсолютного пути к файлу базы данных SQLite путем объединения каталога исполняемого файла с именем файла базы данных.
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), path)

	// Проверка существования файла базы данных SQLite.
	_, err = os.Stat(dbFile)

	// Инициализация флага установки на false.
	var install bool
	if err != nil {
		fmt.Println(err)
		install = true
	}
	return path, install
}

// CheckPort извлекает номер порта из переменной окружения "TODO_PORT".
// Если "TODO_PORT" не установлен, по умолчанию используется "7540".
//
// Возвращает:
// Номер порта в виде строки.
func CheckPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
