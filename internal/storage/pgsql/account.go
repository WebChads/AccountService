package storage

import (
	"database/sql"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (a *AccountRepository) Insert() error {
	// res, err := a.db.Exec("SELECT * FROM accounts")
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return err
	// 	}

	// 	return err
	// }

	// // Get number of affected rows
	// rows, _ := res.RowsAffected()
	// _ = rows

	return nil
}

func (a *AccountRepository) Update() error {
	return nil
}
