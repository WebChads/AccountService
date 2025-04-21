package usecase

import (
	"database/sql"

	"github.com/WebChads/AccountService/internal/models/entities"
	storage "github.com/WebChads/AccountService/internal/storage/pgsql"
)

type AccountRepository interface {
	Insert(*entities.Account) error
}

// All service repositories
type Repositories struct {
	Account AccountRepository
	// ...
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Account: storage.NewAccountRepository(db),
		// ...
	}
}