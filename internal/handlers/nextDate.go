package handlers

import (
	"fmt"
	"go_final_project/internal/utils"
	"net/http"
	"time"
)

// NextDate - это обработчик HTTP-запросов, который вычисляет следующую дату на основе указанных параметров.
// Он ожидает текущую дату и время в параметре "now", целевую дату в параметре "date",
// и частоту повторения в параметре "repeat".
//
// Функция использует предоставленную текущую дату и время для определения следующего вхождения целевой даты.
// Параметр "repeat" может принимать одно из следующих значений: "daily", "weekly", "monthly", или "yearly".
//
// Если входные параметры недействительны или при вычислении возникает ошибка, функция возвращает ответ HTTP 400 Bad Request.
// В противном случае она устанавливает заголовок "Content-Type" в "text/plain" и выводит результат в формате "%s\n".
func NextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\nневерный формат now", err), http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if date == "" || repeat == "" {
		http.Error(w, fmt.Sprintf("%s\nневерный формат date или repeat", err), http.StatusBadRequest)
		return
	}
	next, err := utils.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\n", next)
}
