package service_test

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mystaline/paperid-test/internal/entity"
	"github.com/mystaline/paperid-test/internal/repository"
)

var (
	_ repository.TransactionManager      = stubTxManager{}
	_ repository.UserRepositor           = (*mockUserRepo)(nil)
	_ repository.WalletRepositor         = (*mockWalletRepo)(nil)
	_ repository.TransactionLogRepositor = (*mockTransactionLogRepo)(nil)
)

type stubTxManager struct{}

func (stubTxManager) WithinTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return fn(nil)
}

type mockUserRepo struct {
	getOneUserByID func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error)
}

func (m *mockUserRepo) GetOneUserByID(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
	return m.getOneUserByID(ctx, db, id)
}

type mockWalletRepo struct {
	getWalletByUserID      func(ctx context.Context, db repository.DBTX, userId int64) (*entity.Wallet, error)
	debitBalanceWithReturn func(ctx context.Context, db repository.DBTX, id int64, decrement int64) (*entity.Wallet, error)
}

func (m *mockWalletRepo) GetWalletByUserID(
	ctx context.Context,
	db repository.DBTX,
	userId int64,
) (*entity.Wallet, error) {
	return m.getWalletByUserID(ctx, db, userId)
}

func (m *mockWalletRepo) DebitBalanceWithReturn(
	ctx context.Context,
	db repository.DBTX,
	id int64,
	decrement int64,
) (*entity.Wallet, error) {
	return m.debitBalanceWithReturn(ctx, db, id, decrement)
}

type mockTransactionLogRepo struct {
	insertLog func(ctx context.Context, db repository.DBTX, log entity.TransactionLog) error
}

func (m *mockTransactionLogRepo) InsertLog(ctx context.Context, db repository.DBTX, log entity.TransactionLog) error {
	return m.insertLog(ctx, db, log)
}
