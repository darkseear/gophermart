package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/darkseear/go-musthave/internal/config"
	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/service"
)

type OrderHandler struct {
	orderServices *service.Order
	cfg           *config.Config
}

func NewOrderHandler(orderServices *service.Order, cfg *config.Config) *OrderHandler {
	return &OrderHandler{orderServices: orderServices, cfg: cfg}
}

func (h *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromRequest(r, h.cfg)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	orderNumber := strings.TrimSpace(string(body))
	if orderNumber == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = h.orderServices.UserUploadsOrder(r.Context(), models.Order{Number: orderNumber, UserID: userID, Status: models.Registered})
	if err != nil {
		if errors.Is(err, service.ErrOrderAlreadyExists) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, service.ErrOrderDoesNotExistToAnotherUser) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrInvalidOrder) {
			logger.Log.Error("error upload")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromRequest(r, h.cfg)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderServices.UserGetOrder(r.Context(), userID)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
