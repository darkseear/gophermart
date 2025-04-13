package handlers

import (
	"github.com/darkseear/go-musthave/internal/config"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/processor"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/service"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Router *chi.Mux
	cfg    *config.Config
	store  *repository.Loyalty
}

func Routers(cfg *config.Config, store *repository.Loyalty, auth *service.Auth, processor *processor.Order) *Router {
	r := Router{
		Router: chi.NewRouter(),
		cfg:    cfg,
		store:  store,
	}

	userService := service.NewUser(store)
	userHandler := NewUsersHandler(userService, auth)
	orderService := service.NewOrder(store, processor)
	orderHandler := NewOrderHandler(orderService, r.cfg)
	balanceService := service.NewBalance(store)
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
