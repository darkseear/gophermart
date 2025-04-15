package service

import (
	"context"
	"errors"

	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/utils"
)

type Balance struct {
	store repository.LoyaltyRepository
}

func NewBalance(store repository.LoyaltyRepository) *Balance {
	return &Balance{store: store}
}

func (b *Balance) UserGetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	return b.store.GetBalance(ctx, userID)
}

func (b *Balance) UserGetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return b.store.GetWithdrawals(ctx, userID)
}

func (b *Balance) UserWithdrawn(ctx context.Context, userID int, orderNumber string, amount float64) error {
	if !utils.ValidLuhn(orderNumber) {
		return errors.New("invalid order number")
	}
	if amount <= 0 {
		return errors.New("negative amount")
	}

	err := b.store.CreateWithdrawal(ctx, userID, orderNumber, amount)
	if err != nil {
		if err.Error() == "sql: transaction has already been committed or rolled back" {
			return nil
		}
		if errors.Is(err, errors.New("insufficient funds")) {
			return errors.New("insufficient funds")
		}
		if errors.Is(err, errors.New("failed to create withdrawal")) {
			return errors.New("failed to create withdrawal")
		}
		return err
	}
	return nil
}
