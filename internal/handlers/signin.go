package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) Authentication(w http.ResponseWriter, r *http.Request) {
	// Получаем пароль из переменной окружения
	password := os.Getenv("TODO_PASSWORD")
	var request struct {
		Pass string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	var hashedPass string
	if password == "" {
		hashString := sha256.Sum256([]byte(""))
		hashedPass = hex.EncodeToString(hashString[:])
	} else if password != request.Pass {
		err := errors.New("password is incorrect")
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	if len(password) > 0 {
		hashString := sha256.Sum256([]byte(request.Pass))
		hashedPass = hex.EncodeToString(hashString[:])
	}
	claims := jwt.MapClaims{
		"hashedPass": hashedPass,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// получаем подписанный токен
	signedToken, err := jwtToken.SignedString([]byte(password))
	if err != nil {
		err = fmt.Errorf("failed to sign jwt: %s", err)
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = fmt.Fprintf(w, `{"token": "%s"}`, signedToken)
	if err != nil {
		h.logger.Error(err)
	}
}
