package repository

import (
	"context"

	"github.com/mystaline/paperid-test/internal/entity"
)

var _ TransactionLogRepositor = (*TransactionLogRepository)(nil)

type TransactionLogRepository struct{}

type TransactionLogRepositor interface {
	InsertLog(ctx context.Context, db DBTX, log entity.TransactionLog) error
}

func NewTransactionLogRepository() TransactionLogRepositor {
	return &TransactionLogRepository{}
}

func (r *TransactionLogRepository) InsertLog(
	ctx context.Context,
	db DBTX,
	log entity.TransactionLog,
) error {
	rawQuery := `INSERT INTO transaction_logs (id, user_id, user_full_name, action, amount) VALUES
    ($1, $2, $3, $4, $5)`
	_, err := db.Exec(ctx, rawQuery, log.ID, log.UserID, log.UserFullName, log.Action, log.Amount)
	if err != nil {
		return err
	}

	return nil
}
