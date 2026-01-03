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

type TransactionFunc func(ctx context.Context) error
