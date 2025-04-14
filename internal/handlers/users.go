package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/service"
	"go.uber.org/zap"
)

type UsersHandler struct {
	userService *service.User
	auth        *service.Auth
}

func NewUsersHandler(userService *service.User, auth *service.Auth) *UsersHandler {
	return &UsersHandler{userService: userService, auth: auth}
}

func (uh *UsersHandler) UserRegistration(w http.ResponseWriter, r *http.Request) {

	var userInput models.UserInput
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		logger.Log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := uh.userService.UserRegistration(r.Context(), userInput.Login, userInput.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			logger.Log.Error("User already exists", zap.Error(err))
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else if errors.Is(err, service.ErrUserInvalidPassword) {
			logger.Log.Error("Invalid password", zap.Error(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			logger.Log.Error("Failed to register user", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := uh.auth.GenerateToken(user.ID)
	if err != nil {
		logger.Log.Error("Failed to generate token", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", token)
	logger.Log.Info("User registered and auth", zap.Int("userID", user.ID))
	w.WriteHeader(http.StatusOK)
}

func (uh *UsersHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var userInput models.UserInput
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		logger.Log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := uh.userService.UserLogin(r.Context(), userInput.Login, userInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := uh.auth.GenerateToken(user.ID)
	if err != nil {
		logger.Log.Error("Failed to generate token", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", token)
	logger.Log.Info("User logged in and auth", zap.Int("userID", user.ID))
	w.WriteHeader(http.StatusOK)
}
