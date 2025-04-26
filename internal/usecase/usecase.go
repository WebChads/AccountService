package usecase

import (
	"context"

	"github.com/WebChads/AccountService/internal/models/dtos"
	storage "github.com/WebChads/AccountService/internal/storage/pgsql"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
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
