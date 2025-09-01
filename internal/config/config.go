package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server Server
		Mongo  MongoConfig
		JWT    JWTConfig
	}
	Server struct {
		Port string `mapstructure:"port"`
	}

	MongoConfig struct {
		URI    string
		DBName string `mapstructure:"dbName"`
	}

	JWTConfig struct {
		AccessTokenTTL  string `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL string `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}
)

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found, using system environment variables")
	}

	if err := parseConfigFile("./configs"); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if err := setFromEnv(&cfg); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to set environment variables: %w", err)
	}

	return &cfg, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")
	viper.SetConfigType("yml")

	return viper.ReadInConfig()
}

func setFromEnv(cfg *Config) error {
	cfg.Mongo.URI = os.Getenv("MONGODB_URI")
	cfg.JWT.SigningKey = os.Getenv("SIGNING_KEY")

	if cfg.Mongo.URI == "" {
		return errors.New("MONGODB_URI environment variable is required")
	}
	if cfg.JWT.SigningKey == "" {
		return errors.New("SIGNING_KEY environment variable is required")
	}

	return nil
}
