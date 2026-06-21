package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mystaline/paperid-test/internal/entity"
	"github.com/mystaline/paperid-test/internal/repository"
	"github.com/mystaline/paperid-test/internal/service"
	"github.com/mystaline/paperid-test/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestInvoke_Success(t *testing.T) {
	trnasactionLogRepo := &mockTransactionLogRepo{
		insertLog: func(ctx context.Context, db repository.DBTX, log entity.TransactionLog) error {
			assert.Equal(t, "Disburse", log.Action)
			assert.Equal(t, int64(2000), log.Amount)
			assert.Equal(t, int64(1), log.UserID)
			assert.Equal(t, "User Test", log.UserFullName)
			return nil
		},
	}
	walletRepo := &mockWalletRepo{
		getWalletByUserID: func(ctx context.Context, db repository.DBTX, userId int64) (*entity.Wallet, error) {
			return &entity.Wallet{
				ID:        int64(1),
				UserID:    int64(1),
				Balance:   int64(10000),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
		debitBalanceWithReturn: func(ctx context.Context, db repository.DBTX, id, decrement int64) (*entity.Wallet, error) {
			assert.Equal(t, int64(1), id)
			assert.Equal(t, int64(2000), decrement)
			return &entity.Wallet{
				ID:        int64(1),
				UserID:    int64(1),
				Balance:   int64(8000),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return &entity.User{
				ID:        int64(1),
				FullName:  "User Test",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		trnasactionLogRepo,
		userRepo,
		walletRepo,
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "1", res.ID)
	assert.Equal(t, "1", res.UserID)
	assert.Equal(t, int64(8000), res.RemainingBalance)
}

func TestInvoke_UserNotFound(t *testing.T) {
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return nil, pgx.ErrNoRows
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		&mockTransactionLogRepo{},
		userRepo,
		&mockWalletRepo{},
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.Nil(t, res)
	var httpErr *helper.Error
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestInvoke_UserLookupError(t *testing.T) {
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return nil, errors.New("connection reset")
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		&mockTransactionLogRepo{},
		userRepo,
		&mockWalletRepo{},
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.Nil(t, res)
	var httpErr *helper.Error
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
}

func TestInvoke_WalletNotFound(t *testing.T) {
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return &entity.User{ID: int64(1), FullName: "User Test"}, nil
		},
	}
	walletRepo := &mockWalletRepo{
		getWalletByUserID: func(ctx context.Context, db repository.DBTX, userId int64) (*entity.Wallet, error) {
			return nil, pgx.ErrNoRows
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		&mockTransactionLogRepo{},
		userRepo,
		walletRepo,
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.Nil(t, res)
	var httpErr *helper.Error
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestInvoke_WalletLookupError(t *testing.T) {
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return &entity.User{ID: int64(1), FullName: "User Test"}, nil
		},
	}
	walletRepo := &mockWalletRepo{
		getWalletByUserID: func(ctx context.Context, db repository.DBTX, userId int64) (*entity.Wallet, error) {
			return nil, errors.New("connection reset")
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		&mockTransactionLogRepo{},
		userRepo,
		walletRepo,
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.Nil(t, res)
	var httpErr *helper.Error
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
}

func TestInvoke_InsufficientBalanceOnPreCheck(t *testing.T) {
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return &entity.User{ID: int64(1), FullName: "User Test"}, nil
		},
	}
	walletRepo := &mockWalletRepo{
		getWalletByUserID: func(ctx context.Context, db repository.DBTX, userId int64) (*entity.Wallet, error) {
			return &entity.Wallet{ID: int64(1), UserID: int64(1), Balance: int64(500)}, nil
		},
		debitBalanceWithReturn: func(ctx context.Context, db repository.DBTX, id, decrement int64) (*entity.Wallet, error) {
			t.Fatal("debit must not be called when balance is insufficient")
			return nil, nil
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		&mockTransactionLogRepo{},
		userRepo,
		walletRepo,
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.Nil(t, res)
	var httpErr *helper.Error
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusUnprocessableEntity, httpErr.Code)
}

func TestInvoke_InsufficientBalanceOnDebit(t *testing.T) {
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return &entity.User{ID: int64(1), FullName: "User Test"}, nil
		},
	}
	walletRepo := &mockWalletRepo{
		getWalletByUserID: func(ctx context.Context, db repository.DBTX, userId int64) (*entity.Wallet, error) {
			return &entity.Wallet{ID: int64(1), UserID: int64(1), Balance: int64(10000)}, nil
		},
		debitBalanceWithReturn: func(ctx context.Context, db repository.DBTX, id, decrement int64) (*entity.Wallet, error) {
			return nil, pgx.ErrNoRows
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		&mockTransactionLogRepo{},
		userRepo,
		walletRepo,
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.Nil(t, res)
	var httpErr *helper.Error
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusUnprocessableEntity, httpErr.Code)
}

func TestInvoke_DebitError(t *testing.T) {
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return &entity.User{ID: int64(1), FullName: "User Test"}, nil
		},
	}
	walletRepo := &mockWalletRepo{
		getWalletByUserID: func(ctx context.Context, db repository.DBTX, userId int64) (*entity.Wallet, error) {
			return &entity.Wallet{ID: int64(1), UserID: int64(1), Balance: int64(10000)}, nil
		},
		debitBalanceWithReturn: func(ctx context.Context, db repository.DBTX, id, decrement int64) (*entity.Wallet, error) {
			return nil, errors.New("deadlock detected")
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		&mockTransactionLogRepo{},
		userRepo,
		walletRepo,
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "deadlock detected")
	var httpErr *helper.Error
	assert.False(t, errors.As(err, &httpErr))
}

func TestInvoke_InsertLogError(t *testing.T) {
	transactionLogRepo := &mockTransactionLogRepo{
		insertLog: func(ctx context.Context, db repository.DBTX, log entity.TransactionLog) error {
			return errors.New("insert log failed")
		},
	}
	userRepo := &mockUserRepo{
		getOneUserByID: func(ctx context.Context, db repository.DBTX, id int64) (*entity.User, error) {
			return &entity.User{ID: int64(1), FullName: "User Test"}, nil
		},
	}
	walletRepo := &mockWalletRepo{
		getWalletByUserID: func(ctx context.Context, db repository.DBTX, userId int64) (*entity.Wallet, error) {
			return &entity.Wallet{ID: int64(1), UserID: int64(1), Balance: int64(10000)}, nil
		},
		debitBalanceWithReturn: func(ctx context.Context, db repository.DBTX, id, decrement int64) (*entity.Wallet, error) {
			return &entity.Wallet{ID: int64(1), UserID: int64(1), Balance: int64(8000)}, nil
		},
	}

	svc := service.NewDisbursementService(
		stubTxManager{},
		transactionLogRepo,
		userRepo,
		walletRepo,
	)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: int64(2000), UserID: int64(1)})

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "insert log failed")
}
