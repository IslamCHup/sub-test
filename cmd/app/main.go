package main

import (
	"os"
	"test-junior-go/internal/config"
	postgres "test-junior-go/internal/db"
	"test-junior-go/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	logger := logger.InitLog(cfg.LogLevel)

	_, err := postgres.InitDB(cfg.DB, logger)
	if err != nil {
		logger.Error("failed to initialize db", "error", err)
		return
	}

	r := gin.Default()

	logger.Info("application started")

	if err := r.Run(":" + os.Getenv("APP_PORT")); err != nil {
		logger.Error("", "err", err)
	}
}
