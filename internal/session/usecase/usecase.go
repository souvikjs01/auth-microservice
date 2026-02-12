package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/internal/session"
)

type sessionUC struct {
	sessionRepo session.SessionRepository
	cfg         *config.Config
}

func NewSessionUseCase(sessionRepo session.SessionRepository, cfg *config.Config) session.SessionUseCase {
	return &sessionUC{
		sessionRepo: sessionRepo,
		cfg:         cfg,
	}
}

func (u *sessionUC) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.CreateSession")
	defer span.Finish()

	return u.sessionRepo.CreateSession(ctx, session, expire)
}

func (u *sessionUC) GetSessionID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.GetSessionID")
	defer span.Finish()

	return u.sessionRepo.GetSessionByID(ctx, sessionID)
}

func (u *sessionUC) DeleteByID(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.DeleteByID")
	defer span.Finish()

	return u.sessionRepo.DeleteByID(ctx, sessionID)
}
