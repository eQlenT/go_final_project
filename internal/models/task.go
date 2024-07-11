package models

import (
	"fmt"
	"go_final_project/internal/utils"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Repeat  string `json:"repeat"`
	Comment string `json:"comment"`
}

// CompleteRequest вычисляет следующую дату для указанной задачи на основе предоставленной даты и правила повторения.
// Если поле даты не указано или оно пустое, используется текущая дата.
// Если дата меньше текущей даты, есть два варианта:
// - Если правило повторения не указано или оно пустое, используется текущая дата.
// - Если указано правило повторения, необходимо вычислить и сохранить в базе данных следующую дату, которая больше текущей даты.
//
// Параметры:
// r: Структура Task, содержащая идентификатор, название, дату и правило повторения задачи.
//
// Возвращает:
// Строка, представляющая следующую дату для задачи, или ошибку, если вычисление даты завершается с ошибкой.
// Ошибка будет равна nil, если вычисление даты выполнено успешно.
func (t *Task) CompleteRequest() (string, error) {
	var nextDate string
	now := time.Now().Format("20060102")
	timeNow, _ := time.Parse("20060102", now)
	// Если поле date не указано или содержит пустую строку, берётся сегодняшнее число.
	if t.Date == "" || len(t.Date) == 0 {
		t.Date = time.Now().Format("20060102")
		return t.Date, nil
	}
	nextDate = t.Date
	// Если дата меньше сегодняшнего числа, есть два варианта:
	if date, err := time.Parse("20060102", t.Date); err == nil && date.Before(timeNow) {
		// если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число;
		if t.Repeat == "" || len(t.Repeat) == 0 {
			nextDate = time.Now().Format("20060102")
			// при указанном правиле повторения вам нужно вычислить и записать в таблицу дату выполнения,
			// которая будет больше сегодняшнего числа
		} else if date.Equal(timeNow) {
			return t.Date, nil
		} else {
			nextDate, err = utils.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return "", err
			}
		}
	} else if err == nil && time.Now().Before(date) {
		return t.Date, nil
	}

	return nextDate, nil
}

// CheckRequest - это функция, которая проверяет входные данные для задачи.
// Она проверяет поля ID, названия, даты и повторения в предоставленной задаче.
//
// Параметры:
// r - структура Task, содержащая идентификатор, название, дату и правило повторения задачи.
//
// Возвращает:
// Ошибку, если проверка не пройдена, или nil, если проверка пройдена успешно.
//
// Ошибка может возникать в следующих случаях:
// - если поле ID не является числом, возвращается ошибка "не удается разобрать ID";
// - если поле названия пустое или содержит только пробелы, возвращается ошибка "не указано название задачи";
// - если поле даты не пустое и не соответствует формату "20060102", возвращается ошибка "неверный формат даты";
// - если поле повторения не пустое и не соответствует определенным правилам, возвращается ошибка "неверный формат повторения".
func (t Task) CheckTask() error {
	if t.ID != "" || len(t.ID) != 0 {
		_, err := strconv.Atoi(t.ID)
		if err != nil {
			err = fmt.Errorf("can not parse ID")
			return err
		}
	}
	if len(t.Title) == 0 || t.Title == "" || t.Title == " " {
		return fmt.Errorf("не указано название задачи")
	}
	if t.Date != "" {
		if _, err := time.Parse("20060102", t.Date); err != nil {
			return fmt.Errorf("неверный формат даты %s", t.Date)
		}
	}
	if len(t.Repeat) != 0 || t.Repeat != "" {
		repeatSlc := strings.Split(t.Repeat, " ")
		rule := repeatSlc[0]
		if rule == "y" || rule == "d" || rule == "w" || rule == "m" {
			if len(repeatSlc) > 3 || rule == "y" && len(repeatSlc) > 1 || rule == "d" && len(repeatSlc) == 1 || rule == "d" && len(repeatSlc) > 2 || rule == "w" && len(repeatSlc) != 2 {
				return fmt.Errorf("неверный формат repeat")
			}
		} else {
			return fmt.Errorf("неверный формат repeat")
		}
	}
	return nil
}

func (t Task) CheckDate() (string, error) {
	var date string
	now, err := time.Parse("20060102", time.Now().Format("20060102"))
	if err != nil {
		return "", err
	}
	tmpDate, err := time.Parse("20060102", t.Date)
	if err != nil {
		return "", err
	}
	if t.Date == time.Now().Format("20060102") || now.Before(tmpDate) {
		date = t.Date
	}
	if tmpDate.Before(time.Now()) && !now.Equal(tmpDate) {
		return "", fmt.Errorf("date is less than today's date")
	}
	return date, nil
}
