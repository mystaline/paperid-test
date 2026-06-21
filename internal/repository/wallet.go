package repository

import (
	"context"

	"github.com/mystaline/paperid-test/internal/entity"
)

var _ WalletRepositor = (*WalletRepository)(nil)

type WalletRepository struct{}

type WalletRepositor interface {
	GetByUserID(ctx context.Context, db DBTX, userId int64) (*entity.Wallet, error)
	UpdateOneByID(ctx context.Context, db DBTX, id int64, values entity.Wallet) (*entity.Wallet, error)
}

func NewWalletRepository() WalletRepositor {
	return &WalletRepository{}
}

func (r *WalletRepository) GetByUserID(
	ctx context.Context,
	db DBTX,
	userId int64,
) (*entity.Wallet, error) {
	rawQuery := `SELECT id, user_id, balance FROM wallets WHERE user_id = $1 LIMIT 1`
	rows := db.QueryRow(ctx, rawQuery, userId)

	wallet := entity.Wallet{}
	err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepository) UpdateOneByID(
	ctx context.Context,
	db DBTX,
	id int64,
	values entity.Wallet,
) (*entity.Wallet, error) {
	panic("")
}
