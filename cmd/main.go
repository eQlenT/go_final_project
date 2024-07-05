package main

import (
	"go_final_project/internal/handlers"
	"go_final_project/internal/utils"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	port := utils.CheckPort()
	db := utils.InitDB()
	DBconnection := &handlers.DBConnection{DB: db}
	// Путь к директории веб-файлов
	webDir := "web"

	// Создаем новый экземпляр сервера
	server := &http.Server{
		Addr: ":" + port, // Порт, на котором будет слушаться сервер
	}

	// Настраиваем маршрутизацию для всех файлов в директории web
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.NextDate)
	http.HandleFunc("/api/task", DBconnection.Task)
	http.HandleFunc("/api/tasks", DBconnection.GetTasks)
	http.HandleFunc("/api/task/done", DBconnection.TaskDone)
	// Запускаем сервер
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
