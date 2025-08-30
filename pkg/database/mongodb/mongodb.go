package mongodb

import (
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Initializes a MongoDB connection using the provided configuration parameters.

func NewClient(uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	logger.Info("Connections to MongoDB successfully.")

	return client, nil
}
