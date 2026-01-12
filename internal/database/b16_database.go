package database

import (
	"context"
	"database/sql"
)

type database struct {
	db *sql.DB
}

type databaseTransactionKeyType struct{}

var databaseTransactionKey = databaseTransactionKeyType{}

func NewDatabase(db *sql.DB) Database {
	return &database{db: db}
}

func (db *database) Executor(ctx context.Context) QueryExecutor {
	if tx, ok := ctx.Value(databaseTransactionKey).(*sql.Tx); ok {
		return tx
	}
	return db.db
}

func (db *database) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, databaseTransactionKey, tx)

	defer tx.Rollback()

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
