package models

import "time"

type Withdrawal struct {
	UserID      int       `json:"-" db:"user_id"`
	OrderNumber int       `json:"order" db:"order_number"`
	Sum         float64   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

type ReqWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
