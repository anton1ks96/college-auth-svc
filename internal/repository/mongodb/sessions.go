package mongodb

import (
	"context"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SessionRepository struct {
	cfg *config.Config
	db  *mongo.Database
}

func NewSessionRepository(cfg *config.Config, db *mongo.Database) *SessionRepository {
	return &SessionRepository{
		cfg: cfg,
		db:  db,
	}
}

func (s *SessionRepository) SaveRefreshToken(ctx context.Context, session *domain.RefreshSession) error {

	return nil
}
