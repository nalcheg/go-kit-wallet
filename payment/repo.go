package payment

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

const (
	DefaultOffset = uint64(0)
	DefaultLimit  = uint64(10)
)

type (
	PaymentsResponse struct {
		Total    uint64
		Payments []*Payment
	}
	repo struct {
		db     *sql.DB
		logger log.Logger
	}
)

func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) GetPayments(req *GetPaymentsRequest) (*PaymentsResponse, error) {
	limit := DefaultLimit
	if req.PerPage != nil {
		limit = *req.PerPage
	}
	offset := DefaultOffset
	if req.Page != nil {
		offset = (*req.Page - 1) * limit
	}

	where := ` WHERE 1=1 `
	if req.Direction != nil {
		if *req.Direction == Incoming {
			where += ` AND to_account IS NULL `
		} else if *req.Direction == Outgoing {
			where += ` AND from_account IS NULL `
		}
	}

	query := `SELECT id, amount, account, from_account, to_account FROM payments `
	countQuery := `SELECT COUNT(id) FROM payments `

	var rows *sql.Rows
	var err error
	var total uint64
	if req.AccountID != nil {
		where += ` AND account = $1 `
		if err := repo.db.QueryRow(countQuery+where, req.AccountID).Scan(&total); err != nil {
			return nil, err
		}
		query += where + ` LIMIT $2 OFFSET $3`
		rows, err = repo.db.Query(query, req.AccountID, limit, offset)
	} else {
		if err := repo.db.QueryRow(countQuery + where).Scan(&total); err != nil {
			return nil, err
		}
		query += where + ` LIMIT $1 OFFSET $2`
		rows, err = repo.db.Query(query, limit, offset)
	}
	if err != nil {
		return nil, err
	}

	var resp PaymentsResponse
	resp.Total = total

	for rows.Next() {
		var p Payment
		var amount uint64
		if err := rows.Scan(&p.ID, &amount, &p.Account, &p.FromAccount, &p.ToAccount); err != nil {
			return nil, err
		}

		if p.ToAccount == nil {
			p.Direction = "incoming"
		} else if p.FromAccount == nil {
			p.Direction = "outgoing"
		}

		p.Amount = float64(amount) / 100

		resp.Payments = append(resp.Payments, &p)
	}

	return &resp, nil
}

func (repo *repo) DoPayment(fromAccountID, toAccountID string, amount uint64) error {
	tx, err := repo.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}

	initialState, err := repo.getInitialState(tx, fromAccountID, toAccountID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if initialState.fromCurrency != initialState.toCurrency {
		return NoCurrenciesExchange
	}

	if amount > initialState.fromBalance-amount {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return NeedGreaterBalanceError
	}

	if err := repo.insertPayments(tx, Outgoing, amount, fromAccountID, toAccountID); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := repo.insertPayments(tx, Incoming, amount, toAccountID, fromAccountID); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := repo.updateAccountsBalances(tx, fromAccountID, initialState.fromBalance-amount); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := repo.updateAccountsBalances(tx, toAccountID, initialState.toBalance+amount); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

type Direction bool

const (
	Outgoing Direction = true
	Incoming Direction = false
)

type initialState struct {
	fromBalance, toBalance   uint64
	fromCurrency, toCurrency string
}

func (repo *repo) insertPayments(tx *sql.Tx, direction Direction, amount uint64, accountID, secondAccountID string) error {
	var dstAcc string
	if direction == Incoming {
		dstAcc = "from_account"
	} else if direction == Outgoing {
		dstAcc = "to_account"
	}

	_, err := tx.Exec(
		`INSERT INTO payments (id, amount, account, `+dstAcc+`) VALUES ($1, $2, $3, $4)`,
		uuid.New(), amount, accountID, secondAccountID,
	)

	return err
}

func (repo *repo) updateAccountsBalances(tx *sql.Tx, account string, balance uint64) error {
	query := `UPDATE accounts SET balance = $1 WHERE id = $2`
	_, err := tx.Exec(
		query,
		balance, account,
	)

	return err
}

func (repo *repo) getInitialState(tx *sql.Tx, fromAccountID, toAccountID string) (*initialState, error) {
	query := `
		SELECT balance, currency, 0 sort_order FROM accounts WHERE id = $1
		UNION ALL
		SELECT balance, currency, 1 sort_order FROM accounts WHERE id = $2
		ORDER BY sort_order
	`
	q, err := tx.Query(query, fromAccountID, toAccountID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	var fromBalance, toBalance, order uint64
	var fromCurrency, toCurrency string
	q.Next()
	if err := q.Scan(&fromBalance, &fromCurrency, &order); err != nil {
		return nil, err
	}
	q.Next()
	if err := q.Scan(&toBalance, &toCurrency, &order); err != nil {
		return nil, err
	}

	return &initialState{
		fromBalance:  fromBalance,
		toBalance:    toBalance,
		fromCurrency: fromCurrency,
		toCurrency:   toCurrency,
	}, nil
}
