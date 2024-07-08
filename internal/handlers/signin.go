package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/internal/utils"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Пусть пароль хранится в переменной окружения TODO_PASSWORD.
// Если это значение не пустое — нужно запросить пароль.
// При этом ваши API-запросы тоже должны проверять, аутентифицирован ли пользователь или нет.
// Разберём пошагово реализацию аутентификации.
func Authentication(w http.ResponseWriter, r *http.Request) {
	// Получаем пароль из переменной окружения
	password := os.Getenv("TODO_PASSWORD")
	var request struct {
		Pass string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendErr(w, err, http.StatusBadRequest)
		return
	}

	if password == "" {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if password != request.Pass {
		err := errors.New("password is incorrect")
		utils.SendErr(w, err, http.StatusBadRequest)
		return
	}
	hashString := sha256.Sum256([]byte(request.Pass))
	// так как result — массив байт, а EncodeToString принимает слайс, преобразуем массив в слайс при помощи [:]
	hashedPass := hex.EncodeToString(hashString[:])
	claims := jwt.MapClaims{
		"hashedPass": hashedPass,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// получаем подписанный токен
	signedToken, err := jwtToken.SignedString([]byte(password))
	if err != nil {
		err = fmt.Errorf("failed to sign jwt: %s", err)
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, `{"token": "%s"}`, signedToken)
}
