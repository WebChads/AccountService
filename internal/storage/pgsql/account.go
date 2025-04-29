package storage

import (
	"context"
	"errors"
	"time"

	"github.com/WebChads/AccountService/internal/models/dtos"
	"github.com/google/uuid"
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

// Database inner structure
type Account struct {
	UserId     uuid.UUID `db:"user_id"`
	Firstname  string    `db:"firstname"`
	Surname    string    `db:"surname"`
	Patronymic string    `db:"patronymic"`
	Gender     string    `db:"gender"`
	Birthdate  time.Time `db:"birthdate"`
}

func newAccount(a dtos.CreateAccountRequest) Account {
	// "02-01-2006" is just format string
	birthdate, _ := time.Parse("02-01-2006", a.Birthdate)

	return Account{
		UserId:     a.UserId,
		Firstname:  a.Firstname,
		Surname:    a.Surname,
		Patronymic: a.Patronymic,
		Gender:     a.Gender,
		Birthdate:  birthdate,
	}
}

func (r *AccountRepository) checkExistence(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE user_id = $1)`

	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, errors.New("error checking account existence: " + err.Error())
	}

	return exists, nil
}

func (r *AccountRepository) Select(ctx context.Context, userId uuid.UUID) (*dtos.GetAccountResponse, error) {
	query := `
		SELECT user_id, firstname, surname, patronymic, gender, birthdate
	 	FROM accounts WHERE user_id = :user_id
	`

	params := map[string]any{"user_id": userId}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, errors.New("failed to execute query: " + err.Error())
	}
	defer rows.Close()

	// Check if there are any rows
	if !rows.Next() {
		return nil, errors.New("no account with such id")
	}

	// Process row
	var account Account
	err = rows.StructScan(&account)
	if err != nil {
		return nil, errors.New("failed to get account: " + err.Error())
	}

	response := &dtos.GetAccountResponse{
		UserId:     account.UserId,
		Firstname:  account.Firstname,
		Surname:    account.Surname,
		Patronymic: account.Patronymic,
		Gender:     account.Gender,
		Birthdate:  account.Birthdate,
	}

	return response, nil
}

func (r *AccountRepository) Insert(ctx context.Context, a dtos.CreateAccountRequest) error {
	// Check if user account already exists
	exists, err := r.checkExistence(ctx, a.UserId)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("account with this id already exists")
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
			user_id,
			firstname,
			surname,
			patronymic,
			gender,
			birthdate
	    ) VALUES (:user_id, :firstname, :surname, :patronymic, :gender, :birthdate)
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
