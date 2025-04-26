package dtos

type CreateAccountRequest struct {
	PhoneNumber string
	Firstname   string `json:"firstname" validate:"required"`
	Surname     string `json:"surname" validate:"required"`
	Patronymic  string `json:"patronymic"`
	Gender      string `json:"gender" validate:"required"`
	BirthDate   string `json:"birthdate" validate:"required"`
}
