package usecase

import (
	"context"
	"log/slog"

	"github.com/WebChads/AccountService/internal/models/dtos"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
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

func (a *AccountUsecase) Create(ctx context.Context, account dtos.CreateAccountRequest) error {
	err := a.repository.Insert(ctx, account)
	if err != nil {
		a.logger.Error("create account", slogerr.Error(err))
		return err
	}

	return nil
}
