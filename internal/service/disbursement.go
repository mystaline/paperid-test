package service

import (
	"strconv"

	"github.com/mystaline/paperid-test/internal/dto"
	"github.com/mystaline/paperid-test/internal/repository"
)

type DisbursementParam struct {
	Amount int64
	UserID int64
}

type DisbursementService struct {
	WalletRepository repository.WalletRepositor
}

func NewDisbursementService(walletRepository repository.WalletRepositor) *DisbursementService {
	return &DisbursementService{
		WalletRepository: walletRepository,
	}
}

func (s *DisbursementService) Invoke(param DisbursementParam) (*dto.DisbursementResponse, error) {
	// This is just a placeholder implementation.
	response := &dto.DisbursementResponse{
		ID:     "1",
		Amount: param.Amount,
		UserID: strconv.FormatInt(param.UserID, 10),
	}

	return response, nil
}
