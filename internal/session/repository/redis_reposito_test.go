package repository

import (
	"context"
	"log"

	"testing"

	"github.com/alicebob/miniredis"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/internal/session"
	"github.com/stretchr/testify/require"
)

func SetupRedis() session.SessionRepository {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("Failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	sessionRepo := NewSessionRepository(client, nil)
	return sessionRepo
}

func TestCreateSession(t *testing.T) {
	t.Parallel()

	sessionRepo := SetupRedis()

	t.Run("CreateSession", func(t *testing.T) {
		sessionUUID := uuid.New()
		session := &models.Session{
			SessionID: sessionUUID.String(),
			UserID:    sessionUUID.String(),
		}

		s, err := sessionRepo.CreateSession(context.Background(), session, 10)
		require.NoError(t, err)
		require.NotEqual(t, s, "")

	})
}
