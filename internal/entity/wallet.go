package entity

type Wallet struct {
	ID      int64 `json:"id"`
	UserID  int64 `json:"userId"`
	Balance int64 `json:"balance"`
}
