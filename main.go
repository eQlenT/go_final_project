package main

import (
	"go_final_project/cmd/handlers"
	"go_final_project/cmd/utils"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	port := utils.CheckPort()
	utils.CheckDB()
	// Путь к директории веб-файлов
	webDir := "./web"

	// Создаем новый экземпляр сервера
	server := &http.Server{
		Addr: ":" + port, // Порт, на котором будет слушаться сервер
	}

	// Настраиваем маршрутизацию для всех файлов в директории web
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.NextDate)

	// Запускаем сервер
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
