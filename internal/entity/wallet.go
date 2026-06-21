package entity

import "time"

type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
