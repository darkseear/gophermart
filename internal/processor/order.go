package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/darkseear/go-musthave/internal/accrual"
	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/repository"
	"go.uber.org/zap"
)

type Order struct {
	ordersChan    *Chan
	accrualClient *accrual.Client
	store         repository.LoyaltyRepository
}

type Chan struct {
	orders chan string
	done   chan struct{}
}

func NewOrder(accrualClient *accrual.Client, store repository.LoyaltyRepository) *Order {
	return &Order{
		accrualClient: accrualClient,
		ordersChan: &Chan{
			orders: make(chan string, 100),
			done:   make(chan struct{}),
		},
		store: store,
	}
}

func (o *Order) Start(ctx context.Context) {
	go o.ProcessOrders(ctx)
}

func (o *Order) Stop() {
	close(o.ordersChan.done)
}

func (o *Order) OrderCheck(ctx context.Context, orderNumber string) error {
	accrual, err := o.accrualClient.GetAccrual(orderNumber)
	if err != nil {
		return err
	}
	if accrual == nil {
		return nil
	}

	err = o.store.UpdateOrderStatus(ctx, orderNumber, models.Status(accrual.Status), accrual.Accrual)
	if err != nil {
		return fmt.Errorf("faild to update order status: %w", err)
	}
	if models.Status(accrual.Status) == models.Processed && accrual.Accrual > 0 {
		order, err := o.store.GetOrder(ctx, orderNumber)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}
		if err := o.store.UpdateBalance(ctx, order.UserID, accrual.Accrual); err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}
	}
	return nil
}

func (o *Order) ProcessOrders(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	await := make(map[string]time.Time)
	for {
		select {
		case <-ctx.Done():
			return
		case <-o.ordersChan.done:
			return
		case orderNumber := <-o.ordersChan.orders:
			await[orderNumber] = time.Now()
		case <-ticker.C:
			if len(await) == 0 {
				continue
			}
			for orderNumber, lastCheck := range await {
				if time.Since(lastCheck) < 5*time.Second {
					continue
				}
				if err := o.OrderCheck(ctx, orderNumber); err != nil {
					if _, ok := err.(*accrual.RateLimitError); ok {
						time.Sleep(time.Second)
						continue
					}
					logger.Log.Error("faild to check order", zap.String("order", orderNumber), zap.Error(err))
					continue
				}
				order, err := o.store.GetOrder(ctx, orderNumber)
				if err != nil {
					logger.Log.Error("faild to check order status", zap.String("order", orderNumber), zap.Error(err))
					continue
				}
				if order.Status == models.Processed || order.Status == models.Invalid {
					delete(await, orderNumber)
				} else {
					await[orderNumber] = time.Now()
				}
			}
		}
	}
}

func (o *Order) AddOrder(ctx context.Context, number string) {
	select {
	case o.ordersChan.orders <- number:
		logger.Log.Info("order added to processing queue", zap.String("order", number))
	default:
		logger.Log.Warn("order queue is full, unable to add order", zap.String("order", number))
	}
}
