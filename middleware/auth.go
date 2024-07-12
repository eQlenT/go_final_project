package middleware

import (
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
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			fmt.Println(cookie.Value)
			if err == nil {
				jwt = cookie.Value
			}
			// здесь код для валидации и проверки JWT-токена
			valid := validateJWT(jwt, pass)

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func validateJWT(jwtToken, secret string) bool {
	key := []byte(secret)
	fmt.Println(key)
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != "HS256" {
			return nil, jwt.ErrSignatureInvalid
		}
		return key, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}
