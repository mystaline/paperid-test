package repository

import (
	"context"

	"github.com/mystaline/paperid-test/internal/entity"
)

var _ WalletRepositor = (*WalletRepository)(nil)

type WalletRepository struct{}

type WalletRepositor interface {
	GetWalletByUserID(ctx context.Context, db DBTX, userId int64) (*entity.Wallet, error)
	DebitBalanceWithReturn(ctx context.Context, db DBTX, id int64, decrement int64) (*entity.Wallet, error)
}

func NewWalletRepository() WalletRepositor {
	return &WalletRepository{}
}

func (r *WalletRepository) GetWalletByUserID(
	ctx context.Context,
	db DBTX,
	userId int64,
) (*entity.Wallet, error) {
	rawQuery := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE user_id = $1 LIMIT 1`
	rows := db.QueryRow(ctx, rawQuery, userId)

	wallet := entity.Wallet{}
	err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepository) DebitBalanceWithReturn(
	ctx context.Context,
	db DBTX,
	id int64,
	decrement int64,
) (*entity.Wallet, error) {
	rawQuery := `UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE balance >= $3 AND id = $2 RETURNING id, user_id, balance`
	rows := db.QueryRow(ctx, rawQuery, decrement, decrement, id)

	wallet := entity.Wallet{}
	err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
