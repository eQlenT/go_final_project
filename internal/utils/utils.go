package utils

import (
	"database/sql"
	"fmt"
	"go_final_project/internal/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CheckDB() {
	const (
		CreateTableQuery = `CREATE TABLE scheduler (
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		date    CHAR(8) NOT NULL DEFAULT "",
		title   VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT,
		repeat VARCHAR(128) NOT NULL DEFAULT "" 
		);`
	)

	path := os.Getenv("TODO_DBFILE")
	if path == "" {
		path = "scheduler.db"
	}

	// Получение пути к исполняемому файлу
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), path)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		fmt.Println(err)
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX

	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание БД, в случае её отсутвия
	if install {
		if _, err := db.Exec(CreateTableQuery); err != nil {
			fmt.Println(err)
		}

		if _, err = db.Exec(`CREATE INDEX taks_date ON scheduler (date);`); err != nil {
			fmt.Println(err)

		}
	}
}

func CheckPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}

// NextDate вычисляет следующую дату на основе указанного правила повторения и текущей даты.
//
// Параметры:
// now: Текущая дата и время.
// date: Строка даты в формате "20060102".
// repeat: Строка правила повторения.
//
// Возвращает:
// Строка, представляющая следующую дату в формате "20060102", или ошибку, если правило повторения неверно или дата имеет неверный формат.
//
// Функция поддерживает следующие правила повторения:
// - "d": Ежедневное повторение. Следующая дата вычисляется путем добавления указанного количества дней к текущей дате.
// - "w": Еженедельное повторение. Задача назначается в указанные дни недели, где 1 — понедельник, 7 — воскресенье.
// - "m": Ежемесячное повторение. Следующая дата вычисляется путем добавления указанного количества месяцев к текущей дате.
// - "y": Ежегодное повторение. Следующая дата вычисляется путем добавления указанного количества лет к текущей дате.
//
// Если текущая дата меньше вычисленной следующей даты, функция возвращает ошибку.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	daysInt, monthsInt := make([]int, 0, 7), make([]int, 0, 7)
	var err error

	repeatSlc := strings.Split(repeat, " ")
	rule := repeatSlc[0]
	if len(repeatSlc) > 3 || rule == "y" && len(repeatSlc) > 1 || rule == "d" && len(repeatSlc) == 1 || rule == "d" && len(repeatSlc) > 2 || rule == "w" && len(repeatSlc) > 2 {
		return "", fmt.Errorf("неверный формат repeat")
	}

	if len(repeatSlc) > 1 {
		days := strings.Split(repeatSlc[1], ",")
		dayToAppend := 0
		for idx, day := range days {
			dayToAppend, err = strconv.Atoi(day)
			daysInt = append(daysInt, dayToAppend)
			if daysInt[idx] > 400 {
				return "", fmt.Errorf("неверный формат days (d>400)")
			}
			if err != nil {
				return "", err
			}
		}
	}
	if len(repeatSlc) > 2 {
		months := strings.Split(repeatSlc[2], ",")
		for idx, month := range months {
			monthToAppend, err := strconv.Atoi(month)
			monthsInt = append(monthsInt, monthToAppend)
			if monthsInt[idx] > 12 || monthsInt[idx] < 1 {
				return "", fmt.Errorf("неверный формат month")
			}
			if err != nil {
				return "", err
			}
		}
	}
	now, err = time.Parse("20060102", now.Format("20060102"))
	if err != nil {
		return "", fmt.Errorf("%s/nневерный формат now", err)
	}
	dateStart, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("%s/nневерный формат date", err)
	}

	var resDate time.Time
	switch rule {
	case "":
		if dateStart.Before(now) {
			resDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			return resDate.Format("20060102"), nil
		} else {
			resDate = dateStart
			return resDate.Format("20060102"), nil
		}
	case "d":
		if daysInt[0] == 1 {
			if dateStart == now {
				resDate = dateStart
				return resDate.Format("20060102"), nil
			}
			resDate = dateStart.AddDate(0, 0, 1)
			for resDate.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)) {
				resDate = resDate.AddDate(0, 0, daysInt[0])
			}
			return resDate.Format("20060102"), nil
		}
		if dateStart.Equal(now) {
			resDate = dateStart.AddDate(0, 0, daysInt[0])
			return resDate.Format("20060102"), nil
		}
		if now.Before(dateStart) {
			resDate = dateStart.AddDate(0, 0, daysInt[0])
		} else {
			resDate = dateStart
		}
		for resDate.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)) {
			resDate = resDate.AddDate(0, 0, daysInt[0])
		}
		return resDate.Format("20060102"), nil
	case "y":
		resDate = dateStart.AddDate(1, 0, 0)
		for resDate.Before(now) {
			resDate = resDate.AddDate(1, 0, 0)
		}
		return resDate.Format("20060102"), nil
	case "w":
		Weekdays := make([]int, 0, 7)
		Weekdays = append(Weekdays, daysInt...)
		wantedWeekdays := make(map[int]bool, 7)
		for _, day := range Weekdays {
			if day < 1 || day > 7 {
				return "", fmt.Errorf("неверный формат weekdays")
			}
			if day == 7 {
				wantedWeekdays[0] = true
				continue
			}
			wantedWeekdays[day] = true
		}
		if dateStart.Before(now) {
			resDate = now
		} else {
			resDate = dateStart
		}
		for {
			resDate = resDate.AddDate(0, 0, 1)
			if wantedWeekdays[int(resDate.Weekday())] {
				if resDate.Before(now) {
					return "", fmt.Errorf("полученная дата меньше текущей даты")
				}
				return resDate.Format("20060102"), nil
			}
		}
	case "m":
		wantedMonthDays := make(map[int]bool, 31)
		wantedMonths := make(map[int]bool, 12)
		for _, month := range monthsInt {
			if month < 1 || month > 12 {
				return "", fmt.Errorf("неверный формат months")
			}
			wantedMonths[month] = true
		}
		for _, day := range daysInt {
			if day > 31 || day < -2 {
				return "", fmt.Errorf("неверный формат monthDays")
			}
			if day > 0 {
				wantedMonthDays[day] = true
			} else {
				//подсчёт дня месяца для отрицательных значений
				tempDay := countMonthDay(wantedMonths, now, dateStart, day, dateStart.Before(now))
				wantedMonthDays[tempDay] = true
			}
		}

		resDate = dateStart
		if len(wantedMonths) == 0 && (dateStart.Before(now) || dateStart == now) {
			for day := range wantedMonthDays {
				if time.Date(now.Year(), now.Month(), day, 0, 0, 0, 0, time.Local).Before(now) {
					wantedMonths[int(now.Month())+1] = true
				}
			}
			if len(wantedMonths) == 0 {
				wantedMonths[int(now.Month())] = true
			}
		}

		if len(wantedMonths) == 0 && now.Before(dateStart) {
			for day := range wantedMonthDays {
				var skip bool
				tempDate := time.Date(dateStart.Year(), dateStart.Month(), day, 0, 0, 0, 0, time.Local)
				if tempDate.Month() > dateStart.Month() {
					tempDate = tempDate.AddDate(0, 0, -1)
					if tempDate.Month() == 3 && dateStart.Month() == 2 {
						tempDate = tempDate.AddDate(0, 0, -1)
					}
					skip = true
				}
				if tempDate.Before(dateStart) {
					wantedMonths[int(dateStart.Month())] = true
					if daysInt[0] == -1 || daysInt[0] == -2 {
						for i := 1; i <= 12; i++ {
							wantedMonths[i] = true
						}
					}
				} else if tempDate.Day() > dateStart.Day() {
					if skip {
						wantedMonths[int(dateStart.Month())+1] = true
					} else {
						wantedMonths[int(dateStart.Month())] = true
					}
				}
			}
			if len(wantedMonths) == 0 {
				wantedMonths[int(dateStart.Month())] = true
				if daysInt[0] == -1 || daysInt[0] == -2 {
					for i := 1; i <= 12; i++ {
						wantedMonths[i] = true
					}
				}
			}
		}

		for {
			resDate = resDate.AddDate(0, 0, 1)
			if wantedMonths[int(resDate.Month())] && wantedMonthDays[int(resDate.Day())] && now.Before(resDate) {
				if resDate.Before(now) {
					return "", fmt.Errorf("полученная дата меньше текущей даты")
				}
				return resDate.Format("20060102"), nil
			}
		}
	default:
		return "", fmt.Errorf("неверный формат repeat")

	}
}

func countMonthDay(wantedMonths map[int]bool, now time.Time, dateStart time.Time, subDays int, isDateStartBeforeNow bool) int {
	newMap := make(map[int]bool, 12)
	for k, v := range wantedMonths {
		newMap[k] = v
	}
	if isDateStartBeforeNow {
		dateStart = now
	}
	var isFeb bool
	var daysInMonth int
	var desiredMonth int
	currentMonth := int(dateStart.Month())
	isFeb = currentMonth == 2
	if isFeb && dateStart.Year()%4 == 0 {
		if dateStart.Year()%100 == 0 {
			if dateStart.Year()%400 == 0 {
				daysInMonth = 29
			} else {
				daysInMonth = 28
			}
		} else {
			daysInMonth = 29
		}
	} else if isFeb {
		daysInMonth = 28
	}
	if len(newMap) == 0 {
		if isFeb {
			if dateStart.Day() <= 27 && daysInMonth == 29 {
				if subDays == -1 {
					return daysInMonth
				} else {
					return daysInMonth - 1
				}
			} else if dateStart.Day() <= 27 && daysInMonth == 28 {
				if subDays == -1 {
					return daysInMonth
				} else {
					return 30
				}
			}
			if daysInMonth == 29 && dateStart.Day() == 29 {
				if subDays == -1 {
					return 31
				} else {
					return 30
				}
			} else if daysInMonth == 29 && dateStart.Day() == 28 {
				if subDays == -1 {
					return 29
				} else {
					return 30
				}
			} else if daysInMonth == 28 && dateStart.Day() == 28 {
				if subDays == -1 {
					return 31
				} else {
					return 30
				}
			} else if daysInMonth == 28 && dateStart.Day() == 27 {
				if subDays == -1 {
					return 28
				} else {
					return 30
				}
			}
		} else if dateStart.Day() < 29 {
			desiredMonth = currentMonth
			newMap[currentMonth] = true
		} else if dateStart.Day() == 30 && dateStart.AddDate(0, 0, 1).Day() == 31 {
			if subDays == -1 {
				desiredMonth = currentMonth
				newMap[currentMonth] = true
			} else {
				desiredMonth = currentMonth + 1
				if currentMonth == 12 {
					desiredMonth = 1
				}
				newMap[desiredMonth] = true
			}
		} else if dateStart.Day() == 30 && dateStart.AddDate(0, 0, 1).Day() == 1 {
			desiredMonth = currentMonth + 1
			if currentMonth == 12 {
				desiredMonth = 1
			}
			newMap[desiredMonth] = true

		} else if dateStart.Day() == 31 {
			desiredMonth = currentMonth + 1
			if currentMonth == 12 {
				desiredMonth = 1
			}
			newMap[desiredMonth] = true
		} else if dateStart.Day() == 29 && dateStart.AddDate(0, 0, 2).Day() == 31 {
			desiredMonth = currentMonth
			newMap[currentMonth] = true
		} else if dateStart.Day() == 29 && dateStart.AddDate(0, 0, 2).Day() == 1 {
			if subDays == -1 {
				desiredMonth = currentMonth
				newMap[currentMonth] = true
			} else {
				desiredMonth = currentMonth + 1
				if currentMonth == 12 {
					desiredMonth = 1
				}
				newMap[desiredMonth] = true
			}
		}
	}
	if dateStart.Month() == time.Month(desiredMonth) {
		for newMap[int(dateStart.Month())] {
			dateStart = dateStart.AddDate(0, 0, 1)
		}
		dateStart = dateStart.AddDate(0, 0, subDays)
		return dateStart.Day()
	}
	for {
		dateStart = dateStart.AddDate(0, 0, 1)
		if newMap[int(dateStart.Month())] {
			dateStart = dateStart.AddDate(0, 1, subDays)
			return dateStart.Day()
		}
	}
}

func CompleteRequest(r models.Task) (string, error) {
	var nextDate string
	// Если поле date не указано или содержит пустую строку, берётся сегодняшнее число.
	if r.Date == "" || len(r.Date) == 0 {
		r.Date = time.Now().Format("20060102")
		return r.Date, nil
	}
	// Если дата меньше сегодняшнего числа, есть два варианта:
	if date, err := time.Parse("20060102", r.Date); err == nil && date.Before(time.Now()) {
		// если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число;
		if r.Repeat == "" || len(r.Repeat) == 0 {
			nextDate = time.Now().Format("20060102")
			// при указанном правиле повторения вам нужно вычислить и записать в таблицу дату выполнения,
			// которая будет больше сегодняшнего числа
		} else {
			nextDate, err = NextDate(time.Now(), r.Date, r.Repeat)
			if err != nil {
				return "", err
			}
		}
	} else if err == nil && time.Now().Before(date) {
		nextDate, err = NextDate(time.Now(), r.Date, r.Repeat)
		if err != nil {
			return "", err
		}
	}
	return nextDate, nil
}

func CheckRequest(r models.Task) error {
	if r.ID != "" || len(r.ID) != 0 {
		_, err := strconv.Atoi(r.ID)
		if err != nil {
			err = fmt.Errorf("can not parse ID")
			return err
		}
	}
	if len(r.Title) == 0 || r.Title == "" || r.Title == " " {
		return fmt.Errorf("не указано название задачи")
	}
	if r.Date != "" {
		if _, err := time.Parse("20060102", r.Date); err != nil {
			return fmt.Errorf("неверный формат даты")
		}
	}
	if len(r.Repeat) != 0 || r.Repeat != "" {
		repeatSlc := strings.Split(r.Repeat, " ")
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

func SendErr(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), status)
}
