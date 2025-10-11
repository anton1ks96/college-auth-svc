package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server  Server
		Limiter LimiterConfig
		Mongo   MongoConfig
		JWT     JWTConfig
		LDAP    LDAPConfig
		App     App
		Tokens  Tokens
	}
	Server struct {
		Port           string
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		MaxHeaderBytes int
	}

	App struct {
		Test bool
	}

	LimiterConfig struct {
		RPS   int
		Burst int
		TTL   time.Duration
	}

	LDAPConfig struct {
		URL string
	}

	MongoConfig struct {
		URI      string
		DBName   string
		CollName string
	}

	JWTConfig struct {
		AccessTokenTTL  string
		RefreshTokenTTL string
		SigningKey      string
	}

	Tokens struct {
		InternalToken string
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
	//cfg.LDAP.BindPassword = os.Getenv("BIND_PASSWORD")
	//cfg.LDAP.BindUsername = os.Getenv("BIND_USERNAME")

	if cfg.Mongo.URI == "" {
		return errors.New("MONGODB_URI environment variable is required")
	}
	if cfg.JWT.SigningKey == "" {
		return errors.New("SIGNING_KEY environment variable is required")
	}
	if cfg.LDAP.URL == "" {
		return errors.New("LDAP_URL environment variable is required")
	}
	cfg.Tokens.InternalToken = os.Getenv("INTERNAL_SERVICE_TOKEN")
	if cfg.Tokens.InternalToken == "" {
		return errors.New("INTERNAL_SERVICE_TOKEN environment variable is required")
	}
	//if cfg.LDAP.BindPassword == "" {
	//	return errors.New("BIND_PASSWORD environment variable is required")
	//}
	//if cfg.LDAP.BindUsername == "" {
	//	return errors.New("BIND_USERNAME environment variable is required")
	//}

	return nil
}
