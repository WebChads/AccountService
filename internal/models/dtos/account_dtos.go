package dtos

import (
	"time"

	"github.com/google/uuid"
)

// CreateAccountRequest represents account creation data
// swagger:model CreateAccountRequest
type CreateAccountRequest struct {
	UserId     uuid.UUID `json:"-" swaggerignore:"true"`
	Firstname  string    `json:"firstname" validate:"required" example:"Иван"`
	Surname    string    `json:"surname" validate:"required" example:"Иванов"`
	Patronymic string    `json:"patronymic" example:"Иванович"`
	Gender     string    `json:"gender" validate:"required,min=1,max=1" example:"M"`
	Birthdate  string    `json:"birthdate" validate:"required" example:"1990-01-01"`
}

// GetAccountResponse represents account data
// swagger:model GetAccountResponse
type GetAccountResponse struct {
	UserId     uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Firstname  string    `json:"firstname" example:"Иван"`
	Surname    string    `json:"surname" example:"Иванов"`
	Patronymic string    `json:"patronymic" example:"Иванович"`
	Gender     string    `json:"gender" example:"male"`
	Age        int       `json:"age" example:"33"`
	Birthdate  time.Time `json:"birthdate" example:"1990-01-01T00:00:00Z"`
}
