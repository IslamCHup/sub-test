package main

import (
	"test-junior-go/internal/config"
	"test-junior-go/internal/db"
	"test-junior-go/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	logger := logger.InitLog(cfg.LogLevel)

	// Инициализация подключения к БД
	_, err := db.InitDB(cfg.DB, logger)
	if err != nil {
		logger.Error("failed to initialize db", "error", err)
		return
	}

	logger.Info("application started")
}
