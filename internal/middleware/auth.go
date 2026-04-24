package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/config"
	"backend/pkg/utils"
)

type contextKey string

const UserContextKey contextKey = "userClaims"

func RequireAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Unauthorized: invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims, err := utils.ValidateToken(tokenString, cfg.JWTSecret)
			if err != nil {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*utils.JWTClaims)
		if !ok || claims.Role != "admin" {
			http.Error(w, "Forbidden: requires admin privileges", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
