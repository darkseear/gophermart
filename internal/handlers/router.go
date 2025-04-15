package handlers

import (
	"github.com/darkseear/go-musthave/internal/config"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/service"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Router *chi.Mux
	cfg    *config.Config
	store  *repository.Loyalty
}

func Routers(
	cfg *config.Config, store *repository.Loyalty, auth *service.Auth, userService *service.User,
	orderService *service.Order, balanceService *service.Balance) *Router {
	r := Router{
		Router: chi.NewRouter(),
		cfg:    cfg,
		store:  store,
	}

	userHandler := NewUsersHandler(userService, auth)
	orderHandler := NewOrderHandler(orderService, r.cfg)
	balanceHandler := NewBalanceHandler(balanceService, r.cfg)

	r.Router.Post("/api/user/register", userHandler.UserRegistration) //регистрация пользователя
	r.Router.Post("/api/user/login", userHandler.UserLogin)           //аутентификация пользователя

	r.Router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(auth)) //middleware для аутентификации пользователя
		r.Post("/api/user/orders", middleware.HeaderMiddleware(orderHandler.UploadOrder))
		r.Get("/api/user/orders", middleware.HeaderMiddleware(orderHandler.GetOrders))

		r.Get("/api/user/balance", middleware.HeaderMiddleware(balanceHandler.UserGetBalance))
		r.Post("/api/user/balance/withdraw", balanceHandler.UserWithdrawBalance)
		r.Get("/api/user/withdrawals", middleware.HeaderMiddleware(balanceHandler.UserGetWithdrawals))
	})

	return &r
}
