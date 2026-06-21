package dto

type DisbursementRequest struct {
	Amount int64  `json:"amount"`
	UserID string `json:"userId"`
}

type DisbursementResponse struct {
	ID               string `json:"id"`
	RemainingBalance int64  `json:"remainingBalance"`
	UserID           string `json:"userId"`
}
