package server

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/internal/user"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	userService "github.com/souvikjs01/auth-microservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (u *userServer) Register(c context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	span, c := opentracing.StartSpanFromContext(c, "user.Register")
	defer span.Finish()

	u.logger.Infof("Get request %s\n", r.String())

	createdUser, err := u.userUC.Register(c, &models.User{
		Email:     r.GetEmail(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Password:  r.GetPassword(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "userUC.Register: %v", err)
	}

	return &userService.RegisterResponse{
		Email:     createdUser.Email,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Uid:       createdUser.UserID,
	}, nil
}
