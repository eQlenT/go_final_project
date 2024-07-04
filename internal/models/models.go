package models

type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Repeat  string `json:"repeat"`
	Comment string `json:"comment"`
}
