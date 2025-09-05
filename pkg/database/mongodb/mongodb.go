package mongodb

import (
	"context"
	"errors"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewClient(cfg *config.Config) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(cfg.Mongo.URI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, errors.New("database connection error")
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connections to MongoDB successfully.")

	return client, nil
}
