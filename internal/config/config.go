package config

import (
	"flag"
	"os"
)

type Config struct {
	Address              string
	URL                  string
	LogLevel             string
	Database             string
	SecretKey            string
	AccrualSystemAddress string
}

func New() *Config {
	var config Config

	flag.StringVar(&config.Address, "a", "localhost:8081", "server url")
	flag.StringVar(&config.URL, "b", "http://localhost:8081", "last url")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.StringVar(&config.Database, "d", "host=localhost user=postgres password=1234567890 dbname=loyalty sslmode=disable", "Database")
	flag.StringVar(&config.SecretKey, "s", "secretkey", "Key for JWT")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:8080", "Accrual System Address")

	flag.Parse()

	if val, state := os.LookupEnv("RUN_ADDRESS"); state {
		config.Address = val
	}
	if val, state := os.LookupEnv("BASE_URL"); state {
		config.URL = val
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
