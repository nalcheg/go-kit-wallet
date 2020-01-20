package dbconf

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/tern/migrate"
)

func Migrate(db *sql.DB) error {
	conn, err := stdlib.AcquireConn(db)
	if err != nil {
		return err
	}
	defer stdlib.ReleaseConn(db, conn)

	migrator, err := migrate.NewMigrator(context.Background(), conn, "schema_version")
	if err != nil {
		return err
	}

	if err := migrator.LoadMigrations("migrations"); err != nil {
		return err
	}

	if err := migrator.Migrate(context.Background()); err != nil {
		return err
	}

	return nil
}
