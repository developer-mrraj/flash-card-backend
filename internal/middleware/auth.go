package middleware

import (
	"context"
	"log"
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
			// Skip auth for OPTIONS requests
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("RequireAuth: missing authorization header for %s %s", r.Method, r.URL.Path)
				http.Error(w, "Unauthorized: missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Printf("RequireAuth: invalid authorization header format: %s", authHeader)
				http.Error(w, "Unauthorized: invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims, err := utils.ValidateToken(tokenString, cfg.JWTSecret)
			if err != nil {
				log.Printf("RequireAuth: token validation failed: %v", err)
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			log.Printf("RequireAuth: successful authentication for user: %s", claims.UserID)
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
