package config

import (
	"log/slog"
	"os"

	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	LogLevel    string `yaml:"log_level"`
	Address     string `yaml:"address"`
	DatabaseURL string `yaml:"database_url"`
	SecretKey   string `yaml:"secret_key"`
}

func NewServerConfig() *ServerConfig {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		slog.Error("CONFIG_PATH environment variable not set!")
		return nil
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("config file does not exists: " + configPath)
		return nil
	}

	var cfg ServerConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("cannot read config: %s", slogerr.Error(err))
		return nil
	}

	return &cfg
}
