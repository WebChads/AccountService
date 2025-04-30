package usecase

import (
	"context"

	"github.com/WebChads/AccountService/internal/models/dtos"
	storage "github.com/WebChads/AccountService/internal/storage/pgsql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	Select(ctx context.Context, userId uuid.UUID) (*dtos.GetAccountResponse, error)
	Insert(ctx context.Context, account dtos.CreateAccountRequest) error
}

// All service repositories
type Repositories struct {
	Account AccountRepository
	// ...
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Account: storage.NewAccountRepository(db),
		// ...
	}
}
