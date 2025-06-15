package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	LogLevel    string `json:"log_level" env:"LOG_LEVEL"`
	Address     string `json:"address" env:"ADDRESS"`
	DatabaseURL string `json:"database_url" env:"DATABASE_URL"`
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}

	// Try to load from file first
	fileErr := loadFromFile(cfg)

	// Then override with env vars
	envErr := loadFromEnv(cfg)

	// If both failed, return combined error
	if fileErr != nil && envErr != nil {
		slog.Error(fmt.Errorf("failed to load config: %w, %w", fileErr, envErr).Error())
		return nil
	}

	// Validate config
	if err := validateConfig(cfg); err != nil {
		slog.Error(fmt.Errorf("failed to load config: %w, %w", fileErr, envErr).Error())
		return nil
	}

	return cfg
}

func loadFromFile(cfg *ServerConfig) error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		root, err := findModuleRoot(workingDir)
		if err != nil {
			return fmt.Errorf("failed to find module root: %w", err)
		}

		configPath = filepath.Join(root, "configs", "appsettings.json")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(file, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	slog.Info("config loaded from file", "path", configPath)
	return nil
}

func loadFromEnv(cfg *ServerConfig) error {
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return fmt.Errorf("failed to read env vars: %w", err)
	}

	// Check if any env var was actually set
	emptyCfg := &ServerConfig{}
	if *cfg == *emptyCfg {
		return errors.New("no env vars found")
	}

	slog.Info("config overridden from environment variables")
	return nil
}

func validateConfig(cfg *ServerConfig) error {
	var missing []string

	if cfg.LogLevel == "" {
		missing = append(missing, "log_level")
	}
	if cfg.Address == "" {
		missing = append(missing, "address")
	}
	if cfg.DatabaseURL == "" {
		missing = append(missing, "database_url")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}

func findModuleRoot(dir string) (string, error) {
	for i := 0; i < 10; i++ {
		if dir == "" || dir == "/" {
			return "", errors.New("reached root directory without finding go.mod")
		}

		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		dir = filepath.Dir(dir)
	}

	return "", errors.New("max iterations reached while searching for go.mod")
}
