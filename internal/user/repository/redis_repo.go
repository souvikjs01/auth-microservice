package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/internal/user"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
)

const (
	basePrefix = "user:"
)

type userRedisRepository struct {
	redisClient *redis.Client
	basePrefix  string
	logger      logger.Logger
}

func NewUserRedisRepository(redisClient *redis.Client, logger logger.Logger) user.UserRedisRepository {
	return &userRedisRepository{
		redisClient: redisClient,
		basePrefix:  basePrefix,
		logger:      logger,
	}
}

func (r *userRedisRepository) GetByIdCtx(ctx context.Context, key string) *models.User {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRedisRepository.GetByIdCtx")
	defer span.Finish()

	userBytes, err := r.redisClient.Get(ctx, r.createKey(key)).Bytes()
	if err != nil {
		r.logger.Error("userRedisRepository.GetByIdCtx.redisget: %v", err)
		return nil
	}

	user := &models.User{}
	if err := json.Unmarshal(userBytes, user); err != nil {
		r.logger.Errorf("userRedisRepository.GetByIdCtx.json.unmarshaling %v", err)
		return nil
	}

	return user
}

func (r *userRedisRepository) SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRedisRepository.SetUserCtx")
	defer span.Finish()

	userBytes, err := json.Marshal(&user)
	if err != nil {
		r.logger.Errorf("userRedisRepository.SetUserCtx.json.marshaling %v", err)
		return
	}

	if err := r.redisClient.Set(ctx, r.createKey(key), userBytes, time.Duration(seconds)*time.Second).Err(); err != nil {
		r.logger.Errorf("userRedisRepository.SetUserCtx.redis.set %v", err)
	}
}

func (r *userRedisRepository) DeleteUserCtx(ctx context.Context, key string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRedisRepository.DeleteUserCtx")
	defer span.Finish()

	if err := r.redisClient.Del(ctx, r.createKey(key)).Err(); err != nil {
		r.logger.Errorf("userRedisRepository.DeleteUserCtx.redis.del %v", err)
	}
}

// helper method:
func (r *userRedisRepository) createKey(key string) string {
	return fmt.Sprintf("%s:%s", r.basePrefix, key)
}
