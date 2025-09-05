package mongodb

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
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

func (s *SessionsRepository) SaveRefreshToken(ctx context.Context, session *domain.RefreshSession) error {
	return nil
}
