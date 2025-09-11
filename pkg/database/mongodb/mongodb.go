package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewClient(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(cfg.Mongo.URI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connections to MongoDB successfully.")

	return client, nil
}
