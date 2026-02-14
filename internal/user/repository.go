package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/souvikjs01/auth-microservice/internal/models"
)

type UserPGRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindBYEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
