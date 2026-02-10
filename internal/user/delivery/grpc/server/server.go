package server

import (
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/user"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	userService "github.com/souvikjs01/auth-microservice/proto"
)

type userServer struct {
	userService.UnimplementedUserServiceServer
	logger logger.Logger
	cfg    *config.Config
	userUC user.UserUseCase
}

func NewAuthServerGrpc(logger logger.Logger, cfg *config.Config, userUC user.UserUseCase) *userServer {
	return &userServer{
		logger: logger,
		cfg:    cfg,
		userUC: userUC,
	}
}
