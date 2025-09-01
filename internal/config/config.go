package config

import (
	"errors"
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
		logger.Fatal(errors.New("failed to load environment file"))
	}

	if err := parseConfigFile("./configs"); err != nil {
		logger.Error(errors.New("failed to parse configuration file"))
		return nil, errors.New("failed to parse configuration file")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error(errors.New("failed to unmarshal configuration"))
		return nil, errors.New("failed to unmarshal configuration")
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")
	viper.SetConfigType("yml")

	return viper.ReadInConfig()
}

func setFromEnv(cfg *Config) {
	cfg.Mongo.URI = os.Getenv("MONGODB_URI")
	cfg.JWT.SigningKey = os.Getenv("SIGNING_KEY")
}
