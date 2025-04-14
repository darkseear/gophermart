package config

import (
	"flag"
	"os"
)

type Config struct {
	Address              string `env:"RUN_ADDRESS" envDefault:":8081"`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"info"`
	Database             string `env:"DATABASE_URI" envDefault:"host=localhost user=postgres password=1234567890 dbname=loyalty sslmode=disable"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
	SecretKey            string `env:"SECRET_KEY" envDefault:"secretkey"`
}

func New() *Config {
	var config Config

	flag.StringVar(&config.Address, "a", "localhost:8081", "server url")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.StringVar(&config.Database, "d", "host=localhost user=postgres password=1234567890 dbname=loyalty sslmode=disable", "Database")
	flag.StringVar(&config.SecretKey, "s", "secretkey", "Key for JWT")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:8080", "Accrual System Address")

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
