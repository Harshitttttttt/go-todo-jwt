package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Harshitttttttt/go-todo-jwt/auth"
	"github.com/google/uuid"
)

// Key type for Context Values
type contextKey string

const (
	// UserID Key is the key for userID in the request context
	UserIDKey contextKey = "userID"
)

// AuthMiddleware checks JWT Token
func AuthMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract Token from Authorization Header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization Header Required", http.StatusUnauthorized)
				return
			}

			// Check Bearer Token Format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization Format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate the token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid Or Expired Token", http.StatusUnauthorized)
				return
			}

			// Extract userID from claims
			userIDStr, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid Token Claims", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				http.Error(w, "Invalid UserID in token", http.StatusUnauthorized)
				return
			}

			// Add user ID to the request context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// Call the next handler with the enhanced context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retreives the user ID from the request context
func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	return userID, ok
}
