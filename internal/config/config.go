package config

import (
	"os"

	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	MongoDB struct {
		URI    string `mapstructure:"uri"`
		DBName string `mapstructure:"dbName"`
	} `mapstructure:"mongodb"`
}

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	if err := parseConfigFile("./configs"); err != nil {
		logger.Error(err)
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error(err)
		return nil, err
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
	cfg.MongoDB.URI = os.Getenv("MONGODB_URI")
}
