package models

type Balance struct {
	Current   float64 `json:"balance" db:"current_balance"`
	Withdrawn float64 `json:"withdrawn" db:"current_withdrawn"`
}
