package mongodb

import (
	"context"
	"fmt"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SessionsRepository struct {
	cfg *config.Config
	db  *mongo.Client
}

func NewSessionsRepository(cfg *config.Config, db *mongo.Client) *SessionsRepository {
	return &SessionsRepository{
		cfg: cfg,
		db:  db,
	}
}

func (s *SessionsRepository) SaveRefreshToken(ctx context.Context, session domain.RefreshSession) error {
	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)

	sendSession := domain.RefreshSession{
		JTI:       session.JTI,
		Username:  session.Username,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}

	_, err := coll.InsertOne(ctx, sendSession)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to insert refersh session to MongoDB: %w", err)
	}

	return nil
}

func (s *SessionsRepository) GetRefreshToken(ctx context.Context, jti string) (*domain.RefreshSession, error) {
	return nil, nil
}

func (s *SessionsRepository) RevokeRefreshToken(ctx context.Context, jti string) error {
	return nil
}
