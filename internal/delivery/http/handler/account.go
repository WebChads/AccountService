package handler

import (
	"net/http"
)

type AccountUsecase interface {
	Create()
	Update()
}

type AccountHandler struct {
	usecase AccountUsecase
}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{}
}

func (a *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	a.usecase.Create()
}

func (a *AccountHandler) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	a.usecase.Update()
}
