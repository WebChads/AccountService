package storage

import (
	"context"
	"time"

	"github.com/WebChads/AccountService/internal/models/entities"
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

func (a *AccountRepository) Insert(account *entities.Account) error {
	// Create transactions up to 100 milliseconds long
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond * 100))
    defer cancel()

    query := `
        INSERT INTO accounts (
            firstname, 
            surname, 
            patronymic, 
            gender, 
            birthdate
        ) VALUES (:firstname, :surname, :patronymic, :gender, :birthdate)
        RETURNING id, created_at, updated_at
    `

    // Using NamedQuery + Get for simpler single-row results
    stmt, err := a.db.PrepareNamedContext(ctx, query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    return stmt.GetContext(ctx, account, account)
}
