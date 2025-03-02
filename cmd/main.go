package main

import (
	"log"
	"os"
	"os/signal"
	"quiz-app-be/internal/config"
	"quiz-app-be/internal/handlers"
	"quiz-app-be/internal/repository"
	"quiz-app-be/internal/service"
	"quiz-app-be/internal/setup/aws"
	"quiz-app-be/internal/setup/db"
	"quiz-app-be/internal/setup/httpServer"
	"syscall"

	"go.uber.org/zap"
)

var logger *zap.Logger
var quit = make(chan os.Signal, 1)

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
}

func main() {
	cfg, err := config.Init()
	if err != nil {
		logger.Sugar().Fatalf("failed to init config: %w", err)
		return
	}
	service.SetSecretKey(cfg.SecretKey)
	logger.Info("Config initialized")
	_, err = aws.Init(cfg.Aws)
	if err != nil {
		logger.Sugar().Fatalf("failed to connect s3: %w", err)
		return
	}
	logger.Info("S3 connected")
	pg, err := db.Init(cfg.DB)
	if err != nil {
		logger.Sugar().Fatalf("failed to connect database: %w", err)
		return
	}
	userService := service.NewUserService(repository.NewUsers(pg))
	_, err = httpServer.Init(cfg.HTTPServer, cfg, handlers.NewHandler(
		userService,
	).Routing)
	if err != nil {
		logger.Sugar().Fatalf("init project httpServer error: %w", err)
		return
	}
	logger.Sugar().Infof("Project server started on: %s://%s:%d", cfg.HTTPServer.Proto, cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Printf("SIGTERM received, stopping application...\n")
}
