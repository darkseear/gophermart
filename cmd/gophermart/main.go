package main

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/darkseear/go-musthave/internal/accrual"
	"github.com/darkseear/go-musthave/internal/config"
	"github.com/darkseear/go-musthave/internal/database"
	"github.com/darkseear/go-musthave/internal/handlers"
	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/processor"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/service"
)

func main() {
	if err := run(); err != nil {
		logger.Log.Error("Start server anormal")
		log.Fatal(err)
	}
}

func run() error {
	config := config.New()
	LogLevel := config.LogLevel
	if err := logger.Initialize(LogLevel); err != nil {
		return err
	}

	//инициализировать дб
	db, err := database.InitDB(config.Database)
	if err != nil {
		logger.Log.Error("Failed to initialize database")
		log.Fatal(err)
	}
	defer db.Close()

	//миграции
	err = database.RunMigrations(db)
	if err != nil {
		logger.Log.Error("Failed to run migrations")
		log.Fatal(err)
	}

	auth := service.NewAuth(config.SecretKey)
	ctx := context.Background()
	store := repository.NewLoyalty(db, ctx)
	accrualClient := accrual.NewClient(config.AccrualSystemAddress)

	orderProcessor := processor.NewOrder(accrualClient, store)
	go orderProcessor.Start(ctx)

	userService := service.NewUser(store)
	orderService := service.NewOrder(store, orderProcessor)
	balanceService := service.NewBalance(store)

	r := handlers.Routers(config, store, auth, userService, orderService, balanceService)

	logger.Log.Info("Running server", zap.String("address", config.Address))
	return http.ListenAndServe(config.Address, r.Router)
}
