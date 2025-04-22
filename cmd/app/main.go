package main

import (
	"log/slog"
	"os"

	"github.com/WebChads/AccountService/internal/config"
	server "github.com/WebChads/AccountService/internal/delivery/http"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	prettylogger "github.com/WebChads/AccountService/pkg/pretty_logger"
)

func main() {
	// Init config
	config := config.NewServerConfig()
	if config == nil {
		return
	}

	// Init logger
	logger := setupLogger(config.LogLevel)

	// Init database
	db, err := server.NewDB(config.DatabaseURL)
	if err != nil {
		logger.Error("failed to create database", slogerr.Error(err))
		return
	}

	// Configure server
	router := server.InitRouter(config, logger, db)
	srv := server.NewServer(router, config.Address)

	// Run server
	logger.Info("server started", "address", config.Address)
	srv.ListenAndServe()
}

const (
	envLocal = "local"
	envStage = "stage"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		handler := prettylogger.NewPrettyHandler(os.Stdout)
		log = slog.New(handler)
	case envStage:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
