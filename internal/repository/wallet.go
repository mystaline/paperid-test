package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mystaline/paperid-test/internal/entity"
)

var _ WalletRepositor = (*WalletRepository)(nil)

type WalletRepository struct {
	db *pgxpool.Pool
}

type WalletRepositor interface {
	GetByUserID(userId int64) entity.Wallet
	UpdateOneByID(id int64, values entity.Wallet) entity.Wallet
}

func NewWalletRepository(pool *pgxpool.Pool) WalletRepositor {
	return &WalletRepository{
		db: pool,
	}
}

func (r *WalletRepository) GetByUserID(userId int64) entity.Wallet {
	panic("")
}

func (r *WalletRepository) UpdateOneByID(id int64, values entity.Wallet) entity.Wallet {
	panic("")
}
