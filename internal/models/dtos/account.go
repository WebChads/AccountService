package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	UserId     uuid.UUID
	Firstname  string `json:"firstname" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic"`
	Gender     string `json:"gender" validate:"required"`
	Birthdate  string `json:"birthdate" validate:"required"`
}

type GetAccountResponse struct {
	UserId     uuid.UUID
	Firstname  string    `json:"firstname"`
	Surname    string    `json:"surname"`
	Patronymic string    `json:"patronymic"`
	Gender     string    `json:"gender"`
	Age        int       `json:"age"`
	Birthdate  time.Time `json:"birthdate"`
}
