package storage

import (
	"database/sql"

	"github.com/WebChads/AccountService/internal/models/entities"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (a *AccountRepository) Insert(account *entities.Account) error {
	return nil
}
