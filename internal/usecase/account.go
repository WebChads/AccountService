package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/WebChads/AccountService/internal/models/dtos"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	"github.com/google/uuid"
)

type AccountUsecase struct {
	logger     *slog.Logger
	repository AccountRepository
}

func NewAccountUsecase(r AccountRepository, l *slog.Logger) *AccountUsecase {
	return &AccountUsecase{
		logger:     l,
		repository: r,
	}
}

func (a *AccountUsecase) Get(ctx context.Context, userId string) (*dtos.GetAccountResponse, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		a.logger.Error("user id parsing error", slogerr.Error(err))
		return nil, err
	}

	account, err := a.repository.Select(ctx, id)
	if err != nil {
		a.logger.Error("get account", slogerr.Error(err))
		return nil, err
	}

	// Calculate user age (hardcode for now)
	account.Age = (time.Now().Year() - account.Birthdate.Year())
	if (time.Now().Month() < account.Birthdate.Month()) {
		account.Age--;
	} else if (time.Now().Month() == account.Birthdate.Month()) {
		if (time.Now().Day() < account.Birthdate.Day()) {
			account.Age--;
		}
	}

	return account, nil
}

func (a *AccountUsecase) Create(ctx context.Context, req dtos.CreateAccountRequest) error {
	err := a.repository.Insert(ctx, req)
	if err != nil {
		a.logger.Error("create account", slogerr.Error(err))
		return err
	}

	return nil
}
