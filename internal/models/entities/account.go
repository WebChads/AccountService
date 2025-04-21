package entities

import "github.com/WebChads/AccountService/internal/models/dtos"

type Account struct {
	firstname  string
	surname    string
	patronymic string
	gender     string
	birthdate  string
}

// Constructor
func NewAccountEntity(req dtos.CreateAccountRequest) *Account {
	return &Account{
		firstname: req.Firstname,
		surname: req.Surname,
		patronymic: req.Patronymic,
		gender: req.Gender,
		birthdate: req.BirthDate,
	}
}

