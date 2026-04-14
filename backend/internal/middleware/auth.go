package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"cbt/backend/internal/services"
)

// ContextKey is a type for custom context keys
type ContextKey string

const UserContextKey ContextKey = "user_id"

// RequireAuth verifies JWT token from Authorization header and adds user info to context.
// This middleware should be used to protect authenticated routes.
func RequireAuth(verifyTokenFunc func(string) (*services.Claims, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "missing authorization header"}`, http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				return
			}

			// Extract Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error": "invalid authorization header"}`, http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				return
			}

			token := parts[1]

			// Verify token and get claims
			claims, err := verifyTokenFunc(token)
			if err != nil {
				fmt.Printf("token verification error: %v\n", err)
				http.Error(w, `{"error": "invalid or expired token"}`, http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				return
			}

			// Add user ID to context from claims
			ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
