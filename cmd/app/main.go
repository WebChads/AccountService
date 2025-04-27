package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/WebChads/AccountService/internal/config"
	server "github.com/WebChads/AccountService/internal/delivery/http"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	prettylogger "github.com/WebChads/AccountService/pkg/pretty_logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/WebChads/AccountService/docs"
)

// @title Swagger API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /api/v1
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
	if err := runMigrations(db.DB, logger); err != nil {
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

func runMigrations(db *sql.DB, logger *slog.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("could not create migration driver", slogerr.Error(err))
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations", "postgres", driver,
	)
	if err != nil {
		logger.Error("could not create migration instanse", slogerr.Error(err))
		return err
	}

	// Apply available migrations
	if err := migration.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Warn("could not apply migrations", slogerr.Error(err))
			return nil
		}

		if strings.Contains(err.Error(), "no such table") {
			logger.Warn(
				"first-time migration, creating schema_migrations table",
				slogerr.Warn(err),
			)

        	if err := migration.Force(1); err != nil {
            	return err
        	}

        	return migration.Up()
		}

		logger.Error("could not apply migrations", slogerr.Error(err))
		return err
	}

	return nil
}
