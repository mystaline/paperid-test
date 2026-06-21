package entity

import "time"

type TransactionLog struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"userId"`
	UserFullName string    `json:"userFullName"`
	Action       string    `json:"action"`
	Amount       int64     `json:"amount"`
	CreatedAt    time.Time `json:"createdAt"`
}
