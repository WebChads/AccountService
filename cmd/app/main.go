package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/WebChads/AccountService/docs"
	"github.com/WebChads/AccountService/internal/config"
	server "github.com/WebChads/AccountService/internal/delivery/http"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	"github.com/WebChads/AccountService/internal/storage/pgsql/migrations"
	prettylogger "github.com/WebChads/AccountService/pkg/pretty_logger"
)

// @title Account Service API
// @version 1.0
// @description API for managing user accounts and authentication

// @contact.name API Support
// @contact.url https://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @schemes http
func main() {
	// Init config
	config := config.NewServerConfig()
	if config == nil {
		return
	}

	// Init logger
	logger := setupLogger(config.LogLevel)

	// Create context
	ctx := context.Background()

	// Init database
	db, err := server.NewDB(ctx, config.DatabaseURL)
	if err != nil {
		logger.Error("failed to create database", slogerr.Error(err))
		return
	}
	defer db.Close()

	// Apply migrations
	if err := migrations.RunMigrations(db.DB, logger); err != nil {
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
