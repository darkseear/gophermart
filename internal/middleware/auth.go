package middleware

import (
	"context"
	"net/http"

	"github.com/darkseear/go-musthave/internal/service"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(a *service.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userID, err := a.ValidateToken(token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(token string, secret string) (int, error) {
	a := service.NewAuth(secret)
	userID, err := a.ValidateToken(token)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
