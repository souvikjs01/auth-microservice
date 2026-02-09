package user

import (
	"context"

	"github.com/souvikjs01/auth-microservice/internal/models"
)

type UserUseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
}
