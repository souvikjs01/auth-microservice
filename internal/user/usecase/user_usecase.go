package usecase

import (
	"context"

	"github.com/google/uuid"
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

	return u.userRepo.Create(ctx, user)
}

// find by email
func (u *UserUsecase) FindBYEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.FindBYEmail")
	defer span.Finish()

	return u.userRepo.FindBYEmail(ctx, email)
}

// find by id
func (u *UserUsecase) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.FindByID")
	defer span.Finish()

	return u.userRepo.FindByID(ctx, userID)
}
