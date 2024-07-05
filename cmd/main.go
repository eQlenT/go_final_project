package main

import (
	"database/sql"
	"go_final_project/internal/handlers"
	"go_final_project/internal/utils"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	webDir := "web"
	port := utils.CheckPort()
	path, install := utils.CheckDB()

	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
	}
	DBconnection := &handlers.DBConnection{DB: db}
	if install {
		DBconnection.InitDB()
	}
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
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
