package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/internal/user"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
)

type UserUsecase struct {
	logger   logger.Logger
	userRepo user.UserRepository
}

func NewUserUsecase(logger logger.Logger, userRepo user.UserRepository) *UserUsecase {
	return &UserUsecase{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (u *UserUsecase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.Register")
	defer span.Finish()

	return u.userRepo.Register(ctx, user)
}
