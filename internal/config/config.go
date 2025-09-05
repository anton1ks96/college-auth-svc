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
		LDAP   LDAPConfig
	}
	Server struct {
		Port string `mapstructure:"port"`
	}

	LDAPConfig struct {
		URL          string
		BaseDN       string `mapstructure:"baseDn"`
		BindPassword string
		BindUsername string
	}

	MongoConfig struct {
		URI      string
		DBName   string
		CollName string
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
	cfg.Mongo.DBName = os.Getenv("MONGODB_DBNAME")
	cfg.Mongo.CollName = os.Getenv("MONGODB_CNAME")
	cfg.JWT.SigningKey = os.Getenv("SIGNING_KEY")
	cfg.LDAP.URL = os.Getenv("LDAP_URL")
	cfg.LDAP.BindPassword = os.Getenv("BIND_PASSWORD")
	cfg.LDAP.BindUsername = os.Getenv("BIND_USERNAME")

	if cfg.Mongo.URI == "" {
		return errors.New("MONGODB_URI environment variable is required")
	}
	if cfg.JWT.SigningKey == "" {
		return errors.New("SIGNING_KEY environment variable is required")
	}
	if cfg.LDAP.URL == "" {
		return errors.New("LDAP_URL environment variable is required")
	}
	if cfg.LDAP.BindPassword == "" {
		return errors.New("BIND_PASSWORD environment variable is required")
	}
	if cfg.LDAP.BindUsername == "" {
		return errors.New("BIND_USERNAME environment variable is required")
	}

	return nil
}
