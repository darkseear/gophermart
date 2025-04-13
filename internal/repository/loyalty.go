package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"go.uber.org/zap"
)

type Loyalty struct {
	db  *sql.DB
	ctx context.Context
}

func NewLoyalty(db *sql.DB, ctx context.Context) *Loyalty {
	return &Loyalty{db: db, ctx: ctx}
}

func (l *Loyalty) GreaterUser(ctx context.Context, user models.UserInput) (*models.User, error) {
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id, 
	login, password_hash, created_at`
	userUser := &models.User{}
	err := l.db.QueryRowContext(ctx, query, user.Login, user.Password).Scan(
		&userUser.ID,
		&userUser.Login,
		&userUser.PasswordHash,
		&userUser.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "users_login_key") {
			logger.Log.Error("user already exists", zap.Error(err))
			return nil, errors.New("user already exists")
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("User retrieved successfully", zap.String("login", user.Login))
	return userUser, nil
}

func (l *Loyalty) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password_hash, created_at 
	FROM users 
	WHERE login = $1`
	user := &models.User{}
	err := l.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found", zap.Error(err))
			return nil, err
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("User retrieved successfully", zap.String("login", user.Login))
	return user, nil
}

func (l *Loyalty) UploadOrder(ctx context.Context, order models.Order) error {
	var isOrderExists sql.NullInt64
	query := `SELECT user_id FROM orders WHERE number = $1`
	err := l.db.QueryRowContext(ctx, query, order.Number).Scan(&isOrderExists)
	if err != nil && err != sql.ErrNoRows {
		logger.Log.Error("Failed to check if order exists", zap.Error(err))
		return err
	}
	if err != sql.ErrNoRows {
		if isOrderExists.Valid && isOrderExists.Int64 == int64(order.UserID) {
			return errors.New("order already exists ")
		}
		return errors.New("order does not exist to another user")
	}

	query = `INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)`
	_, err = l.db.ExecContext(ctx, query, order.Number, order.UserID, order.Status)
	if err != nil {
		logger.Log.Error("Failed to insert order", zap.Error(err))
		return err
	}

	logger.Log.Info("Order uploaded successfully", zap.String("orderNumber", order.Number))
	return nil
}

func (l *Loyalty) GetOrder(ctx context.Context, orderNumber string) (*models.Order, error) {
	query := `SELECT number, user_id, status, accrual, uploaded_at 
	FROM orders 
	WHERE number = $1`

	order := &models.Order{}
	err := l.db.QueryRowContext(ctx, query, orderNumber).Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found", zap.Error(err))
			return nil, err
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}
	return order, nil
}

func (l *Loyalty) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	query := `SELECT number, user_id, status, accrual, uploaded_at
	 FROM orders 
	 WHERE user_id = $1
	 ORDER BY uploaded_at DESC`

	rows, err := l.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			logger.Log.Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error occurred during row iteration", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Orders retrieved successfully", zap.Int("userID", userID))
	return orders, nil
}

func (l *Loyalty) UpdateOrderStatus(ctx context.Context, orderNumber string, status models.Status, accrual float64) error {
	query := `UPDATE orders 
	SET status = $1, accrual = $2 
	WHERE number = $3`

	accrualValue := max(accrual, 0)
	resStatus, err := l.db.ExecContext(ctx, query, status, accrualValue, orderNumber)
	if err != nil {
		logger.Log.Error("Failed to update order status", zap.Error(err))
		return err
	}

	rowsAffected, err := resStatus.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}
	logger.Log.Info("Order status updated successfully")
	return nil
}

func (l *Loyalty) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	query := `SELECT current_balance, withdrawn_balance 
	FROM users 
	WHERE id = $1`
	balance := &models.Balance{}
	err := l.db.QueryRowContext(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found", zap.Error(err))
			return &models.Balance{}, err
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Balance retrieved successfully", zap.Int("userID", userID))
	return balance, nil
}

func (l *Loyalty) UpdateBalance(ctx context.Context, userID int, accrual float64) error {
	if ctx.Err() != nil {
		logger.Log.Error("Context canceled before starting transaction", zap.Error(ctx.Err()))
		return ctx.Err()
	}
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	query := `SELECT current_balance FROM users WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, userID).Scan(&currentBalance)
	if err != nil {
		return err
	}

	if accrual < 0 && currentBalance+accrual < 0 {
		return errors.New("insufficient balance")
	}

	query = `UPDATE users SET current_balance = current_balance + $1 WHERE id = $2`
	_, err = l.db.ExecContext(ctx, query, accrual, userID)
	if err != nil {
		logger.Log.Error("Failed to update balance", zap.Error(err))
		return err
	}

	logger.Log.Info("Balance updated successfully", zap.Int("userID", userID), zap.Float64("delta", accrual))
	return tx.Commit()
}

func (l *Loyalty) CreateWithdrawal(ctx context.Context, userID int, orderNumber string, sum float64) error {
	if ctx.Err() != nil {
		logger.Log.Error("Context canceled before starting transaction", zap.Error(ctx.Err()))
		return ctx.Err()
	}
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	query := `SELECT current_balance FROM users WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, userID).Scan(&currentBalance)
	if err != nil {
		return err
	}
	if currentBalance < sum {
		return errors.New("insufficient balance")
	}

	query = `INSERT INTO withdrawals ( order_number, user_id, sum) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, query, orderNumber, userID, sum)
	if err != nil {
		logger.Log.Error("Failed to create withdrawal", zap.Error(err))
		return errors.New("failed to create withdrawal")
	}
	logger.Log.Info("Withdrawal created successfully", zap.Int("userID", userID), zap.String("orderNumber", orderNumber), zap.Float64("sum", sum))
	_, err = tx.ExecContext(ctx, `UPDATE users SET current_balance = current_balance - $1 WHERE id = $2`, sum, userID)
	if err != nil {
		logger.Log.Error("Failed to update user current_balance", zap.Error(err))
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE users SET withdrawn_balance = withdrawn_balance + $1 WHERE id = $2`, sum, userID)
	if err != nil {
		logger.Log.Error("Failed to update user withdrawn_balance", zap.Error(err))
		return err
	}

	return tx.Commit()
}

func (l *Loyalty) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	query := `SELECT order_number, user_id, sum, processed_at
	 FROM withdrawals 
	 WHERE user_id = $1
	 ORDER BY processed_at DESC`
	rows, err := l.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var withdrawals []models.Withdrawal
	for rows.Next() {
		withdrawal := models.Withdrawal{}
		err := rows.Scan(&withdrawal.OrderNumber, &withdrawal.UserID, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			logger.Log.Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error occurred during row iteration", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Withdrawals retrieved successfully", zap.Int("userID", userID))
	return withdrawals, nil
}

func (l *Loyalty) Ping(ctx context.Context) error {
	err := l.db.PingContext(ctx)
	if err != nil {
		logger.Log.Error("Failed to ping database", zap.Error(err))
		return err
	}
	logger.Log.Info("Database ping successful")
	return nil
}

func (l *Loyalty) Close() error {
	err := l.db.Close()
	if err != nil {
		logger.Log.Error("Failed to close database connection", zap.Error(err))
		return err
	}
	logger.Log.Info("Database connection closed successfully")
	return nil
}
