package repository

import (
	"context"
	"fmt"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)

	filter := bson.M{"username": session.Username}
	deleted, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		logger.Error(fmt.Errorf("failed to delete existing sessions for user %s: %w", session.Username, err))
		return err
	}

	if deleted.DeletedCount > 0 {
		logger.Debug(fmt.Sprintf("deleted %d existing sessions for user %s", deleted.DeletedCount, session.Username))
	}

	_, err = coll.InsertOne(ctx, session)
	if err != nil {
		logger.Error(fmt.Errorf("failed to save refresh session for user %s: %w", session.Username, err))
		return err
	}

	logger.Debug(fmt.Sprintf("refresh session saved for user %s", session.Username))
	return nil
}

//func (s *SessionsRepository) GetRefreshToken(ctx context.Context, jti string) (*domain.RefreshSession, error) {
//	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)
//
//	var session domain.RefreshSession
//	filter := bson.M{"jti": jti}
//
//	err := coll.FindOne(ctx, filter).Decode(&session)
//	if err != nil {
//		if errors.Is(err, mongo.ErrNoDocuments) {
//			logger.Debug(fmt.Sprintf("refresh session not found for JTI: %s", jti))
//			return nil, err
//		}
//		logger.Error(fmt.Errorf("failed to retrieve refresh session for JTI %s: %w", jti, err))
//		return nil, err
//	}
//
//	logger.Debug(fmt.Sprintf("refresh session retrieved for user %s", session.Username))
//	return &session, nil
//}

func (s *SessionsRepository) RevokeRefreshToken(ctx context.Context, jti string) error {
	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)
	filter := bson.M{"jti": jti}

	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error(fmt.Errorf("failed to revoke refresh session for JTI %s: %w", jti, err))
		return err
	}

	logger.Debug(fmt.Sprintf("refresh session revoked for JTI: %s", jti))
	return nil
}
