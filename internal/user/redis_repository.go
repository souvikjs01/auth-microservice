package user

import (
	"context"

	"github.com/souvikjs01/auth-microservice/internal/models"
)

type UserRedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) *models.User
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.User)
	DeleteUserCtx(ctx context.Context, key string)
}
