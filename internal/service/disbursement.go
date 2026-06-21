package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/mystaline/paperid-test/internal/dto"
	"github.com/mystaline/paperid-test/internal/entity"
	"github.com/mystaline/paperid-test/internal/repository"
	"github.com/mystaline/paperid-test/pkg/db"
	"github.com/mystaline/paperid-test/pkg/helper"
)

type DisbursementParam struct {
	Amount int64
	UserID int64
}

type DisbursementService struct {
	txManager                repository.TransactionManager
	TransactionLogRepository repository.TransactionLogRepositor
	UserRepository           repository.UserRepositor
	WalletRepository         repository.WalletRepositor
}

func NewDisbursementService(
	txManager repository.TransactionManager,
	transactionLogRepository repository.TransactionLogRepositor,
	userRepository repository.UserRepositor,
	walletRepository repository.WalletRepositor,
) *DisbursementService {
	return &DisbursementService{
		txManager:                txManager,
		TransactionLogRepository: transactionLogRepository,
		UserRepository:           userRepository,
		WalletRepository:         walletRepository,
	}
}

func (s *DisbursementService) Invoke(ctx context.Context, param DisbursementParam) (*dto.DisbursementResponse, error) {
	var response *dto.DisbursementResponse

	err := s.txManager.WithinTx(ctx, func(tx pgx.Tx) error {
		user, err := s.UserRepository.GetOneUserByID(ctx, tx, param.UserID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return helper.ErrInternal("Something went wrong with server")
		}
		if user == nil {
			return helper.ErrNotFound("User doesn't exists")
		}

		wallet, err := s.WalletRepository.GetWalletByUserID(ctx, tx, param.UserID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return helper.ErrInternal("Something went wrong with server")
		}
		if wallet == nil {
			return helper.ErrNotFound("Wallet not found")
		}

		if wallet.Balance < param.Amount {
			return helper.ErrUnprocessableEntity("Insufficient balance")
		}

		updatedWallet, err := s.WalletRepository.DebitBalanceWithReturn(ctx, tx, wallet.ID, param.Amount)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return helper.ErrUnprocessableEntity("Insufficient balance")
			}
			return err
		}

		response = &dto.DisbursementResponse{
			ID:               strconv.FormatInt(updatedWallet.ID, 10),
			RemainingBalance: updatedWallet.Balance,
			UserID:           strconv.FormatInt(updatedWallet.UserID, 10),
		}

		err = s.TransactionLogRepository.InsertLog(ctx, tx, entity.TransactionLog{
			ID:           db.IDNode().Generate().Int64(),
			UserID:       user.ID,
			UserFullName: user.FullName,
			Action:       "Disburse",
			Amount:       param.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
