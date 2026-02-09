package server

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	userService "github.com/souvikjs01/auth-microservice/proto"
)

type userServer struct {
	userService.UnimplementedUserServiceServer
	logger logger.Logger
	cfg    *config.Config
}

func NewAuthServerGrpc(logger logger.Logger, cfg *config.Config) *userServer {
	return &userServer{
		logger: logger,
		cfg:    cfg,
	}
}

func (u *userServer) Register(c context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	span, c := opentracing.StartSpanFromContext(c, "user.Register")
	defer span.Finish()

	u.logger.Infof("Get request %s\n", r.String())

	return &userService.RegisterResponse{
		Email:     r.GetEmail(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Uid:       "1",
	}, nil
}
