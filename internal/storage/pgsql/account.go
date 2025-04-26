package storage

import (
	"context"
	"errors"

	"github.com/WebChads/AccountService/internal/models/dtos"
	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// func checkExistanse(account dtos.CreateAccountRequest) bool {
// 	return false
// }

func (a *AccountRepository) Insert(ctx context.Context, account dtos.CreateAccountRequest) error {
	// TODO: check if user account already exists

	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}
	
	query := `
		INSERT INTO accounts (
			firstname, 
			surname, 
			patronymic, 
			gender, 
			birthdate
        ) VALUES (:firstname, :surname, :patronymic, :gender, :birthdate)
	`
	
	_, err = tx.NamedExecContext(ctx, query, account)
	if err != nil {
        return errors.New("failed to insert account: " + err.Error())
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return errors.New("failed to commit transaction: " + err.Error())
    }

	return nil
}
