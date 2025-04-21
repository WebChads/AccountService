package dtos

type CreateAccountRequest struct {
	Firstname  string `json:"firstname"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Gender     string `json:"gender"`
	BirthDate  string `json:"birthdate"`
}