package auth

import (
	"bytes"
	"context"
	"encoding/json"

	"net/http"
	"strings"

	"io"

	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	AuthServiceUrl string
	client         *http.Client
}

func NewMiddleware(authServiceUrl string) *Middleware {
	return &Middleware{
		AuthServiceUrl: authServiceUrl,
		client:         &http.Client{},
	}
}

type AuthServiceResponseDto struct {
	IsValid bool `json:"is_valid"`
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Извлекаем токен из заголовка
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. Проверяем формат Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			authHeader = parts[1]
		}

		tokenString := authHeader

		// 3. Отправляем запрос в auth-сервис для валидации
		reqBody, _ := json.Marshal(map[string]string{"token": tokenString})
		resp, err := m.client.Post(
			"http://"+m.AuthServiceUrl+"/api/v1/auth/validate-token",
			"application/json",
			bytes.NewBuffer(reqBody),
		)
		if err != nil {
			http.Error(w, "Auth service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var responseDto AuthServiceResponseDto
		err = json.Unmarshal(b, &responseDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if resp.StatusCode != http.StatusOK || !responseDto.IsValid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 4. Парсим JWT (без валидации, так как auth-сервис уже проверил)
		token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// 5. Извлекаем claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// 6. Достаем user_id и user_role
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "Missing user_id in token", http.StatusUnauthorized)
			return
		}

		userRole, ok := claims["user_role"].(string)
		if !ok {
			userRole = "" // default если не указана
		}

		// 7. Добавляем в контекст (исправленная версия)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", userID)
		ctx = context.WithValue(ctx, "user_role", userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
