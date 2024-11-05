package middleware

import (
	"context"
	"net/http"
	"strings"
	"testovoe/internal/db"
	"testovoe/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

// Middleware для проверки токена пользователя
func TokenAuthMiddleware(dbProvider *db.PostgresProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Получаем пользователя по токену
			user, err := dbProvider.GetUserByToken(token)
			if err != nil || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Добавляем пользователя в контекст
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Получение пользователя из контекста
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	return user, ok
}
