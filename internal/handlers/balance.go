package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/darkseear/go-musthave/internal/config"
	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/service"
)

type BalanceHandler struct {
	balanceService *service.Balance
	cfg            *config.Config
}

func NewBalanceHandler(balanceService *service.Balance, cfg *config.Config) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService, cfg: cfg}
}

func (b *BalanceHandler) UserGetBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromRequest(r, b.cfg)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := b.balanceService.UserGetBalance(r.Context(), userID)
	if err != nil {
		logger.Log.Error("balance error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		logger.Log.Error("faild to decode request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (b *BalanceHandler) UserWithdrawBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromRequest(r, b.cfg)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.ReqWithdraw
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = b.balanceService.UserWithdrawn(r.Context(), userID, req.Order, req.Sum)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBalanceInvalidOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, service.ErrBalanceInvalidNegativeAmount):
			w.WriteHeader(http.StatusPaymentRequired)
		case errors.Is(err, service.ErrBalanceInvalidInsufficientBalance):
			w.WriteHeader(http.StatusPaymentRequired)
		default:
			logger.Log.Error("unexpected error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (b *BalanceHandler) UserGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromRequest(r, b.cfg)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := b.balanceService.UserGetWithdrawals(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
