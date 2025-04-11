package models

type Accrual struct {
	Order   string  `json:"order"`
	Status  Status  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
