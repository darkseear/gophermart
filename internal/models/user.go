package models

import "time"

type User struct {
	ID           int       `json:"-" db:"id"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
}

type UserInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
