package server

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/internal/models"
	userService "github.com/souvikjs01/auth-microservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u *userServer) Register(c context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	span, _ := opentracing.StartSpanFromContext(c, "user.Register")
	defer span.Finish()

	u.logger.Infof("Get request %s\n", r.String())

	user, err := u.registerRequestToUserModel(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "registerReqToUserModel: %v", err)
	}

	return &userService.RegisterResponse{
		User: u.userModelToProto(user),
	}, nil
}

func (u *userServer) registerRequestToUserModel(r *userService.RegisterRequest) (*models.User, error) {
	candidate := &models.User{
		Email:     r.GetEmail(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Password:  r.GetPassword(),
		Avatar:    r.GetAvatar(),
		Role:      r.GetRole(),
	}

	if err := candidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return candidate, nil
}

func (u *userServer) userModelToProto(user *models.User) *userService.User {
	userProto := &userService.User{
		Uid:       user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		Avatar:    user.Avatar,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
	return userProto
}
