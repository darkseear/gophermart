package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	logger "github.com/darkseear/go-musthave/internal/logging"
)

func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Log.Error("Ошибка подключения к базе данных")
		return nil, err
	}
	return db, nil
}
