package user

import (
	"context"

	"github.com/souvikjs01/auth-microservice/internal/models"
)

type UserRedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*models.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}
