package service

import (
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/metric"
	"github.com/souvikjs01/auth-microservice/internal/session"
	"github.com/souvikjs01/auth-microservice/internal/user"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	userProtoService "github.com/souvikjs01/auth-microservice/proto"
)

type userService struct {
	userProtoService.UnimplementedUserServiceServer
	logger    logger.Logger
	cfg       *config.Config
	userUC    user.UserUseCase
	sessionUC session.SessionUseCase
	mtr       metric.Metrics
}

func NewAuthServerGrpc(logger logger.Logger, cfg *config.Config, userUC user.UserUseCase, sessionUC session.SessionUseCase, mtr metric.Metrics) *userService {
	return &userService{
		logger:    logger,
		cfg:       cfg,
		userUC:    userUC,
		sessionUC: sessionUC,
		mtr:       mtr,
	}
}
