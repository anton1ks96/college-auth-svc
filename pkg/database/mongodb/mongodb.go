package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	if err := createIndexes(client, cfg); err != nil {
		logger.Error(fmt.Errorf("failed to create indexes: %w", err))
	}

	return client, nil
}

func createIndexes(client *mongo.Client, cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(cfg.Mongo.DBName).Collection(cfg.Mongo.CollName)

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userid", Value: 1}},
			Options: options.Index().SetName("userid_idx"),
		},
		{
			Keys:    bson.D{{Key: "jti", Value: 1}},
			Options: options.Index().SetName("jti_idx").SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().
				SetName("expires_at_idx").
				SetExpireAfterSeconds(0),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	logger.Info("MongoDB indexes created successfully")
	return nil
}
