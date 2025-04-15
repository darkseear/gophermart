package config

import (
	"flag"
	"os"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	LogLevel             string `env:"LOG_LEVEL"`
	Database             string `env:"DATABASE_URI" `
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey            string `env:"SECRET_KEY"`
}

func New() *Config {
	config := Config{
		Address:              "localhost:8081",
		LogLevel:             "info",
		Database:             "host=localhost user=postgres password=1234567890 dbname=loyalty sslmode=disable",
		AccrualSystemAddress: "http://localhost:8080",
		SecretKey:            "secretkey",
	}

	flag.StringVar(&config.Address, "a", config.Address, "server url")
	flag.StringVar(&config.LogLevel, "l", config.LogLevel, "log level")
	flag.StringVar(&config.Database, "d", config.Database, "Database")
	flag.StringVar(&config.SecretKey, "s", config.SecretKey, "Key for JWT")
	flag.StringVar(&config.AccrualSystemAddress, "r", config.AccrualSystemAddress, "Accrual System Address")

	flag.Parse()

	if val, state := os.LookupEnv("RUN_ADDRESS"); state {
		config.Address = val
	}
	if val, state := os.LookupEnv("LOG_LEVEL"); state {
		config.LogLevel = val
	}
	if val, state := os.LookupEnv("DATABASE_URI"); state {
		config.Database = val
	}
	if val, state := os.LookupEnv("SECRET_KEY"); state {
		config.SecretKey = val
	}
	if val, state := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); state {
		config.AccrualSystemAddress = val
	}

	return &config
}
