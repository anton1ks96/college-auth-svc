package config

import (
	"os"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server Server
		Mongo  MongoConfig
		JWT    JWTConfig
		Auth   AuthConfig
	}
	Server struct {
		Port string `mapstructure:"port"`
	}

	MongoConfig struct {
		URI    string
		DBName string `mapstructure:"dbName"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	AuthConfig struct {
		JWT          JWTConfig
		PasswordSalt string
	}
)

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(domain.ErrEnvFileLoadFailed)
	}

	if err := parseConfigFile("./configs"); err != nil {
		logger.Error(domain.ErrConfigParsingFailed)
		return nil, domain.ErrConfigParsingFailed
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error(domain.ErrConfigUnmarshalFailed)
		return nil, domain.ErrConfigUnmarshalFailed
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
	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
}
