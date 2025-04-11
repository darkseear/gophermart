package service

import (
	"context"
	"errors"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/utils"
)

type User struct {
	store repository.LoyaltyRepository
}

func NewUser(store repository.LoyaltyRepository) *User {
	return &User{store: store}
}

func (u *User) UserRegistration(ctx context.Context, login, password string) (*models.User, error) {
	passwordHash := utils.HashPassword(password)
	user, err := u.store.GreaterUser(ctx, models.UserInput{Login: login, Password: passwordHash})
	if err != nil {
		if err.Error() == "user already exists" {
			logger.Log.Error("User already exists")
			return nil, errors.New("user already exists")
		}
		return nil, err
	}
	return user, nil
}

func (u *User) UserLogin(ctx context.Context, login, password string) (*models.User, error) {
	user, err := u.store.GetUserByLogin(ctx, login)
	if err != nil {
		logger.Log.Error("Failed to get user by login")
		return nil, err
	}
	if user.PasswordHash != utils.HashPassword(password) {
		return nil, errors.New("invalid password")
	}
	return user, nil
}
