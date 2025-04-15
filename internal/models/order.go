package models

import "time"

type Status string

const (
	Registered Status = "NEW"
	Peocessing Status = "PROCESSING"
	Invalid    Status = "INVALID"
	Processed  Status = "PROCESSED"
)

type Order struct {
	Number     string    `json:"number" db:"number"`
	UserID     int       `json:"-" db:"user_id"`
	Status     Status    `json:"status" db:"status"`
	Accrual    float64   `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
