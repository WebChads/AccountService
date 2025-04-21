package usecase

import (
	"github.com/WebChads/AccountService/internal/models/dtos"
	"github.com/WebChads/AccountService/internal/models/entities"
)

type AccountUsecase struct {
	repository AccountRepository
}

func NewAccountUsecase(r AccountRepository) *AccountUsecase {
	return &AccountUsecase{
		repository: r,
	}
}

func (a *AccountUsecase) Create(req dtos.CreateAccountRequest) error {
	// get entity from DTO
	account := entities.NewAccountEntity(req)
	_ = account

	a.repository.Insert(account)
	return nil
}
