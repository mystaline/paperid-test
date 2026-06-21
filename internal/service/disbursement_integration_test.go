package service_test

import (
	"cmp"
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mystaline/paperid-test/internal/config"
	"github.com/mystaline/paperid-test/internal/entity"
	"github.com/mystaline/paperid-test/internal/repository"
	"github.com/mystaline/paperid-test/internal/service"
	"github.com/mystaline/paperid-test/pkg/db"
)

const (
	itUserID   = int64(800000001)
	itWalletID = int64(800000002)
)

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg := config.LoadConfig()
	cfg.DBHost = cmp.Or(os.Getenv("TEST_DB_HOST"), "localhost")
	cfg.DBPort = cmp.Or(os.Getenv("TEST_DB_PORT"), "55432")

	pool, err := db.NewPool(cfg)
	if err != nil {
		t.Skipf("skip integration: db unreachable (run `docker compose up`): %v", err)
	}
	return pool
}

func seedWallet(t *testing.T, pool *pgxpool.Pool, balance int64) {
	t.Helper()
	ctx := context.Background()

	_, err := pool.Exec(
		ctx,
		`DELETE FROM transaction_logs WHERE user_id = $1`,
		itUserID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx,
		`INSERT INTO users (id, full_name) VALUES ($1, 'IT User')
		 ON CONFLICT (id) DO UPDATE SET full_name = EXCLUDED.full_name`,
		itUserID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx,
		`INSERT INTO wallets (id, user_id, balance) VALUES ($1, $2, $3)
		 ON CONFLICT (id) DO UPDATE SET balance = EXCLUDED.balance`,
		itWalletID,
		itUserID,
		balance,
	)
	require.NoError(t, err)
}

func walletBalance(t *testing.T, pool *pgxpool.Pool) int64 {
	t.Helper()
	var balance int64

	require.NoError(
		t,
		pool.QueryRow(context.Background(), `SELECT balance FROM wallets WHERE id = $1`, itWalletID).Scan(&balance),
	)

	return balance
}

func realService(
	pool *pgxpool.Pool,
	log repository.TransactionLogRepositor, // mockable to trigger failure, to test rollback
) *service.DisbursementService {
	return service.NewDisbursementService(
		repository.NewTransactionManager(pool),
		log,
		repository.NewUserRepository(),
		repository.NewWalletRepository(),
	)
}

func TestInvoke_Integration_Success(t *testing.T) {
	pool := newPool(t)
	defer pool.Close()
	seedWallet(t, pool, 10000)

	svc := realService(pool, repository.NewTransactionLogRepository())

	res, err := svc.Invoke(
		context.Background(),
		service.DisbursementParam{Amount: 2000, UserID: itUserID},
	)

	require.NoError(t, err)
	assert.Equal(t, int64(8000), res.RemainingBalance)
	assert.Equal(t, int64(8000), walletBalance(t, pool))
}

func TestInvoke_Integration_RollbackOnLogFailure(t *testing.T) {
	pool := newPool(t)
	defer pool.Close()
	seedWallet(t, pool, 10000)

	failingLog := &mockTransactionLogRepo{
		insertLog: func(ctx context.Context, db repository.DBTX, log entity.TransactionLog) error {
			return errors.New("forced log failure")
		},
	}
	svc := realService(pool, failingLog)

	res, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: 2000, UserID: itUserID})

	require.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, int64(10000), walletBalance(t, pool))
}

func TestInvoke_Integration_ConcurrentNoOverdraft(t *testing.T) {
	pool := newPool(t)
	defer pool.Close()
	seedWallet(t, pool, 10000)

	svc := realService(pool, repository.NewTransactionLogRepository())

	const amount, workers = int64(1000), 20
	var wg sync.WaitGroup
	var mu sync.Mutex
	success := 0

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.Invoke(context.Background(), service.DisbursementParam{Amount: amount, UserID: itUserID})
			if err == nil {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	// 10 is expected success loop count for 1000 to used all 10000 balance
	assert.Equal(t, 10, success)
	assert.Equal(t, int64(0), walletBalance(t, pool))
}
