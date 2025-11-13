package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/anton1ks96/college-auth-svc/internal/config"
	"github.com/anton1ks96/college-auth-svc/internal/domain"
	"github.com/anton1ks96/college-auth-svc/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	_, err := coll.InsertOne(ctx, session)
	if err != nil {
		logger.Error(fmt.Errorf("failed to save refresh session for user %s: %w", session.UserID, err))
		return err
	}

	return nil
}

func (s *SessionsRepository) TokenExists(ctx context.Context, jti string) (bool, error) {
	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)

	var result struct {
		UserID string `bson:"userid"`
	}

	filter := bson.M{"jti": jti}

	err := coll.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"userid": 1})).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *SessionsRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)
	filter := bson.M{"userid": userID}

	result, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("revoked %d sessions for user %s", result.DeletedCount, userID))
	return nil
}

func (s *SessionsRepository) RevokeRefreshToken(ctx context.Context, jti string) error {
	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)
	filter := bson.M{"jti": jti}

	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error(fmt.Errorf("failed to revoke refresh session for JTI %s: %w", jti, err))
		return err
	}

	return nil
}

func (s *SessionsRepository) ReplaceRefreshToken(ctx context.Context, oldJTI string, newSession *domain.RefreshSession) error {
	ok, err := s.TokenExists(ctx, oldJTI)
	if err != nil {
		logger.Error(fmt.Errorf("failed to check token existence for JTI %s: %w", oldJTI, err))
		return fmt.Errorf("failed to check token existence: %w", err)
	}

	if !ok {
		logger.Debug(fmt.Sprintf("token %s already does not exist, nothing to replace", oldJTI))
		return fmt.Errorf("old token not found: %s", oldJTI)
	}

	if err := s.RevokeRefreshToken(ctx, oldJTI); err != nil {
		logger.Error(fmt.Errorf("failed to delete old token %s during replacement: %w", oldJTI, err))
		return fmt.Errorf("failed to delete old token: %w", err)
	}

	if err := s.SaveRefreshToken(ctx, newSession); err != nil {
		logger.Error(fmt.Errorf("failed to save new token after deleting old one. old_jti=%s, new_jti=%s, user=%s, error=%w",
			oldJTI, newSession.JTI, newSession.UserID, err))

		return fmt.Errorf("failed to save new token (old token was already deleted): %w", err)
	}

	return nil
}

func (s *SessionsRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)

	var session domain.RefreshSession
	filter := bson.M{"userid": userID}

	err := coll.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user session not found")
		}
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	return &domain.User{
		ID:       session.UserID,
		Username: session.Username,
		Role:     session.Role,
	}, nil
}

func (s *SessionsRepository) GetExtendedUserByID(ctx context.Context, userID string) (*domain.UserExtended, error) {
	coll := s.db.Database(s.cfg.Mongo.DBName).Collection(s.cfg.Mongo.CollName)

	var session domain.RefreshSession
	filter := bson.M{"userid": userID}

	err := coll.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user session not found")
		}
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	return &domain.UserExtended{
		ID:            session.UserID,
		Username:      session.Username,
		Role:          session.Role,
		AcademicGroup: session.AcademicGroup,
		Profile:       session.Profile,
		Subgroup:      session.Subgroup,
		EnglishGroup:  session.EnglishGroup,
	}, nil
}
