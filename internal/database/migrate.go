package database

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	logger "github.com/darkseear/go-musthave/internal/logging"
)

//go:embed migrations/*
var migrations embed.FS

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log.Error("не удалось инициализировать драйвер для миграций")
		return err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		logger.Log.Error("не удалось инициализировать источник миграций")
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		logger.Log.Error("не удалось создать объект миграциий")
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Error("не удалось выполнить миграции")
		return err
	}

	logger.Log.Info("Миграции успешно выполнены")
	return nil
}
