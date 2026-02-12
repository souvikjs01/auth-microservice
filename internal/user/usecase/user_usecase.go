package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

	findEmail, err := u.userRepo.FindBYEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.FindByEmail")
	}

	findEmail.SanitizePassword()

	return findEmail, nil
}

// find by id
func (u *UserUsecase) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.FindByID")
	defer span.Finish()

	return u.userRepo.FindByID(ctx, userID)
}

// login user
func (u *UserUsecase) Login(ctx context.Context, email, password string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.Login")
	defer span.Finish()

	findUser, err := u.userRepo.FindBYEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.FindByEmail")
	}

	if err := findUser.ComparePasswords(password); err != nil {
		return nil, errors.Wrap(err, "user.ComparePasswords")
	}

	return findUser, nil
}
