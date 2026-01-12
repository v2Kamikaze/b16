package database

import (
	"context"
)

type Database interface {
	Executor(ctx context.Context) QueryExecutor
	WithTransaction(ctx context.Context, fn TransactionFunc) error
}

type TransactionFunc func(ctx context.Context) error
