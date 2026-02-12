package session

import (
	"context"

	"github.com/souvikjs01/auth-microservice/internal/models"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
