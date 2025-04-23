package usecase

import (
	"github.com/WebChads/AccountService/internal/models/entities"
	storage "github.com/WebChads/AccountService/internal/storage/pgsql"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	Insert(*entities.Account) error
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
