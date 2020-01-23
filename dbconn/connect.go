package dbconn

import (
	"database/sql"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func Connect(dbDSN string) (*sql.DB, error) {
	connConfig, err := pgx.ParseConfig(dbDSN)
	if err != nil {
		return nil, err
	}

	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
