package api

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"go-api-boilerplate/internal/jwt"
	"go-api-boilerplate/internal/shared"
)

func apiKeyMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			header := r.Header.Get("x-api-key")
			if header == "" {
				shared.WriteError(w, http.StatusUnauthorized, "missing api key")
				return
			}
			if subtle.ConstantTimeCompare([]byte(header), []byte(apiKey)) != 1 {
				shared.WriteError(w, http.StatusUnauthorized, "invalid api key")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func authMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				shared.WriteError(w, http.StatusUnauthorized, "missing token")
				return
			}

			userID, err := jwt.ParseAccessToken(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")), jwtSecret)
			if err != nil {
				shared.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			next.ServeHTTP(w, r.WithContext(shared.WithUserID(r.Context(), userID)))
		})
	}
}

func isPublicPath(path string) bool {
	switch path {
	case "/health", "/api/v1/auth/login", "/api/v1/auth/refresh", "/api/v1/auth/logout", "/api/v1/auth/forgot-password":
		return true
	default:
		return false
	}
}
