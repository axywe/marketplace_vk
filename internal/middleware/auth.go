package middleware

import (
	"context"
	"marketplace/internal/store"
	"marketplace/pkg/jwt"
	"marketplace/util"
	"net/http"
	"strings"
)

func Auth(next http.HandlerFunc, db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			util.SendJSONError(w, r, "Authorization header format must be 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		username, err := jwt.ParseToken(headerParts[1], db.Config.JWTSecret)
		if err != nil {
			util.SendJSONError(w, r, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), util.ContextKey("username"), username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func OptionalAuth(next http.HandlerFunc, db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				username, err := jwt.ParseToken(headerParts[1], db.Config.JWTSecret)
				if err == nil {
					ctx := context.WithValue(r.Context(), util.ContextKey("username"), username)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	}
}
