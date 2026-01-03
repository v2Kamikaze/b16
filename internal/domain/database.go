package domain

import (
	"context"
	"database/sql"
)

type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type Database interface {
	Executor(ctx context.Context) QueryExecutor
	WithTransaction(ctx context.Context, fn TransactionFunc) error
}

type database struct {
	db *sql.DB
}

const DatabaseTransactionKey = "DatabaseTransactionKey"

type TransactionFunc func(ctx context.Context) error

func NewDatabase(db *sql.DB) Database {
	return &database{db: db}
}

func (db *database) Executor(ctx context.Context) QueryExecutor {
	if tx, ok := ctx.Value(DatabaseTransactionKey).(*sql.Tx); ok {
		return tx
	}
	return db.db
}

func (db *database) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, DatabaseTransactionKey, tx)

	defer tx.Rollback()

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
