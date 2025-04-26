package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
}

// AuthMiddleware with JWT validation from Authorization header
func (config *JWTConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check header format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.SecretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Add user information to context
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}
