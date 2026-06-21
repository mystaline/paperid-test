package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TransactionManager interface {
	WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type pgxTransaction struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &pgxTransaction{
		pool: pool,
	}
}

func (at *pgxTransaction) WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := at.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
