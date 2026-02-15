package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/internal/user"
	"github.com/souvikjs01/auth-microservice/pkg/grpc_errors"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
)

const (
	userByIdCashedDuration = 3600
)

type UserUsecase struct {
	logger        logger.Logger
	userPgRepo    user.UserPGRepository
	userRedisRepo user.UserRedisRepository
}

func NewUserUsecase(logger logger.Logger, userRepo user.UserPGRepository, userRedisRepo user.UserRedisRepository) *UserUsecase {
	return &UserUsecase{
		logger:        logger,
		userPgRepo:    userRepo,
		userRedisRepo: userRedisRepo,
	}
}

func (u *UserUsecase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.Register")
	defer span.Finish()

	existedUser, err := u.userPgRepo.FindBYEmail(ctx, user.Email)
	if err == nil && existedUser != nil {
		return nil, grpc_errors.ErrEmailExists
	}

	return u.userPgRepo.Create(ctx, user)
}

// find by email
func (u *UserUsecase) FindBYEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.FindBYEmail")
	defer span.Finish()

	findEmail, err := u.userPgRepo.FindBYEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "userPgRepo.FindByEmail")
	}

	findEmail.SanitizePassword()

	return findEmail, nil
}

// find by id
func (u *UserUsecase) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.FindByID")
	defer span.Finish()

	cashedUser, err := u.userRedisRepo.GetByIdCtx(ctx, userID.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		u.logger.Errorf("userUC.FindByID.Cashed user error: %v", err)
	}

	if cashedUser != nil {
		u.logger.Infof("userUC.FindByID.Cashed user: %v", cashedUser)
		return cashedUser, nil
	}

	foundUser, err := u.userPgRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "userUC.FindByID.userPgRepo.FindByID")
	}

	if err := u.userRedisRepo.SetUserCtx(ctx, userID.String(), userByIdCashedDuration, foundUser); err != nil {
		u.logger.Errorf("userUC.FindByID.SetUserCtx error: %v", err)
	}

	return foundUser, nil
}

// login user
func (u *UserUsecase) Login(ctx context.Context, email, password string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUsecase.Login")
	defer span.Finish()

	findUser, err := u.userPgRepo.FindBYEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "userusecase.login.userPgRepo.FindByEmail")
	}

	if err := findUser.ComparePasswords(password); err != nil {
		return nil, errors.Wrap(err, "user.ComparePasswords")
	}

	return findUser, nil
}
