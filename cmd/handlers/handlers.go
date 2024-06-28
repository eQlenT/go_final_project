package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/cmd/utils"
	"go_final_project/models"
	"net/http"
	"time"
)

func NextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\nневерный формат now", err), http.StatusBadRequest)
		w.Write([]byte(""))
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if date == "" || repeat == "" {
		http.Error(w, fmt.Sprintf("%s\nневерный формат date или repeat", err), http.StatusBadRequest)
		w.Write([]byte(""))
		return
	}
	next, err := utils.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusBadRequest)
		w.Write([]byte(""))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\n", next)
}

// func AddTask(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		err := fmt.Errorf("request method must be POST")
// 		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusBadRequest)
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		fmt.Fprintf(w, `{"error": "%s"}`, err)
// 		return
// 	}
// 	var Request models.Request
// 	// Считываем тело запроса в байты
// 	bodyBytes, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusBadRequest)
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		fmt.Fprintf(w, `{"error": "%s"}`, err)
// 		return
// 	}
// 	defer r.Body.Close()

// 	// Анмаршалл JSON'а в структуру
// 	err = json.Unmarshal(bodyBytes, &Request)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusBadRequest)
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		fmt.Fprintf(w, `{"error": "%s"}`, err)
// 		return
// 	}

// 	err = utils.CheckRequest(Request)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusBadRequest)
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		fmt.Fprintf(w, `{"error": "%s"}`, err)
// 		return
// 	}

// 	db, err := sql.Open("sqlite3", "go_final_project/scheduler.db")
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusInternalServerError)
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		fmt.Fprintf(w, `{"error": "%s"}`, err)
// 		return
// 	}
// 	defer db.Close()

// 	nextDate, err := utils.CompleteRequest(Request)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusBadRequest)
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		fmt.Fprintf(w, `{"error": "%s"}`, err)
// 		return
// 	}
// 	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
// 		sql.Named("date", nextDate),
// 		sql.Named("title", Request.Title),
// 		sql.Named("comment", Request.Comment),
// 		sql.Named("repeat", Request.Repeat))
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusInternalServerError)
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		fmt.Fprintf(w, `{"error": "%s"}`, err)
// 		return
// 	}

//		lastInsertID, err := res.LastInsertId()
//		if err != nil {
//			http.Error(w, fmt.Sprintf("%s\n", err), http.StatusInternalServerError)
//			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//			fmt.Fprintf(w, `{"error": "%s"}`, err)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		fmt.Fprintf(w, `{"id": "%d"}`, lastInsertID)
//	}
func AddTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, "request method must be POST", http.StatusBadRequest)
		return
	}

	var request models.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf("error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	err = utils.CheckRequest(request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf("error validating request: %v", err), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "go_final_project/scheduler.db")
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf("error opening database: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	nextDate, err := utils.CompleteRequest(request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf("error processing request: %v", err), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		nextDate, request.Title, request.Comment, request.Repeat)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf("error inserting into database: %v", err), http.StatusInternalServerError)
		return
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf("error getting last insert ID: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, `{"id": "%d"}`, lastInsertID)
}
