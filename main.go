package main

import (
	"database/sql"
	"go_final_project/internal/handlers"
	"go_final_project/internal/utils"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

// main является точкой входа в приложение. Она инициализирует сервер, настраивает маршрутизацию,
// и запускает прослушивание входящих подключений.
func main() {
	webDir := "./web"                // Каталог, содержащий статические файлы для обслуживания
	port := utils.CheckPort()        // Функция для проверки и возврата номера порта
	path, install := utils.CheckDB() // Функция для проверки и возврата пути к базе данных и флага установки

	// Открываем подключение к базе данных
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем новый экземпляр DBConnection и инициализируем базу данных, если это необходимо
	DBconnection := &handlers.DBConnection{DB: db}
	if install {
		DBconnection.InitDB()
	}

	// Создаем новый экземпляр http.Server с указанным портом
	server := &http.Server{
		Addr: ":" + port, // Порт, на котором сервер будет прослушивать
	}

	// Настраиваем маршрутизацию для обслуживания всех файлов в каталоге web и конечных точек API
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/signin", handlers.Authentication)
	http.HandleFunc("/api/nextdate", handlers.NextDate)
	http.HandleFunc("/api/task", DBconnection.Task)
	http.HandleFunc("/api/tasks", DBconnection.GetTasks)
	http.HandleFunc("/api/task/done", DBconnection.TaskDone)

	// Запускаем сервер и прослушиваем входящие подключения
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
