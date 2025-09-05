package mongodb

import (
	"errors"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Mongo struct {
	cfg *config.Config
}

func (m *Mongo) NewClient() (*mongo.Client, error) {
	opts := options.Client().ApplyURI(m.cfg.Mongo.URI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, errors.New("database connection error")
	}

	db := client.Database(m.cfg.Mongo.DBName)
	db.Collection(m.cfg.Mongo.CollName)

	logger.Info("Connections to MongoDB successfully.")

	return client, nil
}
