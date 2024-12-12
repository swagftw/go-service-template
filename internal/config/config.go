package config

import (
	"github.com/spf13/viper"
	"go-service-template/internal"
	"go-service-template/internal/apperr"
	"log/slog"
	"net/http"
	"os"
)

type Config struct {
	Port        int    `mapstructure:"PORT"`
	Debug       bool   `mapstructure:"DEBUG"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

// InitConfig initializes the config and returns it
// If the ENV is not set, it returns an error
func InitConfig() (*Config, error) {
	env := os.Getenv("ENV")

	if env == "" {
		slog.Error("ENV is not set, value should be one of [development, staging, production]")

		return nil, apperr.New(http.StatusInternalServerError, nil, "ENV is not set, value should be one of [development, staging, production]", apperr.ErrInternalError)
	}

	switch env {
	case string(internal.EnvDevelopment):
		viper.AddConfigPath(".")
		viper.SetConfigFile("development.env")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			slog.Error("failed to read the config file", "err", err)

			return nil, apperr.New(http.StatusInternalServerError, err, "failed to read the config file", apperr.ErrInternalError)
		}
	case string(internal.EnvStaging):
		viper.AutomaticEnv()
	case string(internal.EnvProduction):
		viper.AutomaticEnv()
	}

	cfg := new(Config)

	err := viper.UnmarshalExact(cfg)
	if err != nil {
		slog.Error("failed to unmarshal the config file", "err", err)

		return nil, apperr.New(http.StatusInternalServerError, err, "failed to unmarshal the config file", apperr.ErrInternalError)
	}

	return cfg, nil
}
