package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			hashString := sha256.Sum256([]byte(pass))
			// так как result — массив байт, а EncodeToString принимает слайс, преобразуем массив в слайс при помощи [:]
			hashedPassServer := hex.EncodeToString(hashString[:])
			secret := []byte(pass)
			var signedToken string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				signedToken = cookie.Value
			}
			var valid bool
			jwtToken, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				err = fmt.Errorf("failed to parse token: %s", err)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				http.Error(w, fmt.Sprintf(`{"error": "%s","token": "%s"}`, err.Error(), signedToken), http.StatusUnauthorized)
				return
			}
			valid = jwtToken.Valid
			if valid {
				// приводим поле Claims к типу jwt.MapClaims
				res, ok := jwtToken.Claims.(jwt.MapClaims)

				if !ok {
					err = fmt.Errorf("failed to typecast to jwt.MapCalims")
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusUnauthorized)
					return
				}
				hashedPassRaw := res["hashedPass"]
				hashedPass, ok := hashedPassRaw.(string)
				if !ok {
					err = fmt.Errorf("failed to typecast hashedPass to string")
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusUnauthorized)
					return
				}
				if hashedPassServer != hashedPass {
					http.Error(w, "authentification required, wrong password", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "authentification required, token not valid", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
