package server

import (
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/session"
	"github.com/souvikjs01/auth-microservice/internal/user"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	userService "github.com/souvikjs01/auth-microservice/proto"
)

type userServer struct {
	userService.UnimplementedUserServiceServer
	logger    logger.Logger
	cfg       *config.Config
	userUC    user.UserUseCase
	sessionUC session.SessionUseCase
}

func NewAuthServerGrpc(logger logger.Logger, cfg *config.Config, userUC user.UserUseCase, sessionUC session.SessionUseCase) *userServer {
	return &userServer{
		logger:    logger,
		cfg:       cfg,
		userUC:    userUC,
		sessionUC: sessionUC,
	}
}
