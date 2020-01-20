package account

import (
	"database/sql"

	"github.com/go-kit/kit/log"
)

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) CreateAccount(acc *Account) error {
	if _, err := repo.db.Exec(
		`INSERT INTO accounts (id, balance, currency) VALUES ($1, $2, $3)`,
		acc.ID, acc.Balance, acc.Currency,
	); err != nil {
		return err
	}

	return nil
}

func (repo *repo) GetAccount(id string) (*Account, error) {
	q := repo.db.QueryRow(`SELECT id, balance, currency FROM accounts WHERE id = $1`, id)

	var acc Account
	if err := q.Scan(&acc.ID, &acc.Balance, &acc.Currency); err != nil {
		return nil, err
	}

	return &acc, nil
}
