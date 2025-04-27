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

type Account struct {
	PhoneNumber string `db:"phone_number"`
	Firstname   string `db:"firstname"`
	Surname     string `db:"surname"`
	Patronymic  string `db:"patronymic"`
	Gender      string `db:"gender"`
	Birthdate   string `db:"birthdate"`
}

func newAccount(a dtos.CreateAccountRequest) Account {
	return Account{
		PhoneNumber: a.PhoneNumber,
		Firstname:   a.Firstname,
		Surname:     a.Surname,
		Patronymic:  a.Patronymic,
		Gender:      a.Gender,
		Birthdate:   a.BirthDate,
	}
}

func (r *AccountRepository) checkExistanse(ctx context.Context, phoneNumber string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE phone_number = $1)`

	err := r.db.GetContext(ctx, &exists, query, phoneNumber)
	if err != nil {
		return false, errors.New("error checking account existence: " + err.Error())
	}

	return exists, nil
}

func (r *AccountRepository) Insert(ctx context.Context, a dtos.CreateAccountRequest) error {
	// Check if user account already exists
	exists, err := r.checkExistanse(ctx, a.PhoneNumber)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("account with this phone number already exists")
	}

	account := newAccount(a)

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}
	defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

	query := `
		INSERT INTO accounts (
			phone_number,
			firstname,
			surname,
			patronymic,
			gender,
			birthdate
	    ) VALUES (:phone_number, :firstname, :surname, :patronymic, :gender, :birthdate)
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
