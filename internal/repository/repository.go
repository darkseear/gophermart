package repository

import (
	"context"

	"github.com/darkseear/go-musthave/internal/models"
)

type LoyaltyRepository interface {
	GreaterUser(ctx context.Context, user models.UserInput) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)

	UploadOrder(ctx context.Context, order models.Order) error
	GetOrders(ctx context.Context, userID int) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderNumber string, status models.Status, accrual float64) error
	GetOrder(ctx context.Context, orderNumber string) (*models.Order, error)

	GetBalance(ctx context.Context, userID int) (*models.Balance, error)
	UpdateBalance(ctx context.Context, userID int, delta float64) error

	CreateWithdrawal(ctx context.Context, userID int, orderNumber string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error)

	Ping(ctx context.Context) error
	Close() error
}
