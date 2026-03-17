package middleware

import (
	"auth_service/auth"
	"context"
	"net/http"
	"strings"
)

// Функция проверка токена, middleware который применяется к обработчику
func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Требуется авторизация", 401)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				http.Error(w, "Неверный формат токена", 401)
				return
			}

			claims, err := auth.ValidateToken(token, secret)
			if err != nil {
				http.Error(w, "Недействительный токен", 401)
				return
			}

			if claims.Type != "access" {
				http.Error(w, "Неверный тип токена", 401)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
