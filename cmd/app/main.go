// @title Subscription Aggregator API
// @version 1.0
// @description REST API для управления онлайн подписками пользователей
// @host localhost:8080
// @BasePath /api
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test-junior-go/internal/config"
	"test-junior-go/internal/database"
	"test-junior-go/internal/handler"
	"test-junior-go/internal/logger"
	"test-junior-go/internal/repository"
	"test-junior-go/internal/service"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	cfg := config.Load()

	log := logger.InitLog(cfg.LogLevel)

	db, err := database.InitDB(cfg.DB, log)
	if err != nil {
		log.Error("failed to initialize db", "error", err)
		os.Exit(1)
	}

	repo := repository.NewSubscriptionRepository(db, log)
	svc := service.NewSubscriptionService(repo, log)
	hnd := handler.NewSubscriptionHandler(svc, log)
	router := handler.NewRouter(hnd)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Info("server started", "port", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("server exited")
}
