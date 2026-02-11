package server

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/pkg/grpc_errors"
	"github.com/souvikjs01/auth-microservice/pkg/utils"
	userService "github.com/souvikjs01/auth-microservice/proto"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Register User
func (u *userServer) Register(c context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "user.Register")
	defer span.Finish()

	user, err := u.registerRequestToUserModel(r)
	if err != nil {
		u.logger.Errorf("registerReqToUserModel %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "registerReqToUserModel: %v", err)
	}

	if err := utils.ValidateStruct(ctx, user); err != nil {
		u.logger.Errorf("validation error: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "validation error: %v", err)
	}

	createdUser, err := u.userUC.Register(ctx, user)
	if err != nil {
		u.logger.Errorf("userUC.Register %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "userUC.Register: %v", err)
	}

	return &userService.RegisterResponse{
		User: u.userModelToProto(createdUser),
	}, nil
}

// FindByEmail User
func (u *userServer) FindByEmail(c context.Context, r *userService.FindByEmailRequest) (*userService.FindByEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "user.FindByEmail")
	defer span.Finish()

	email := r.GetEmail()

	if !utils.ValidateEmail(email) {
		u.logger.Errorf("invalid email: %s", email)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(errors.New("ErrInvalid email address")), "invalid email: %s", email)
	}

	user, err := u.userUC.FindBYEmail(ctx, email)
	if err != nil {
		u.logger.Errorf("userUC.FindBYEmail %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "userUC.FindBYEmail: %v", err)
	}

	return &userService.FindByEmailResponse{
		User: u.userModelToProto(user),
	}, nil
}

// find by id
func (u *userServer) FindByID(c context.Context, r *userService.FindByIDRequest) (*userService.FindByIDResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "user.FindByID")
	defer span.Finish()

	userId, err := uuid.Parse(r.GetUid())
	if err != nil {
		u.logger.Errorf("invalid user id: %s", r.GetUid())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "invalid user id: %s", r.GetUid())
	}

	user, err := u.userUC.FindByID(ctx, userId)
	if err != nil {
		u.logger.Errorf("userUC.FindByID %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "userUC.FindByID: %v", err)
	}

	return &userService.FindByIDResponse{
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
		Uid:       user.UserID.String(),
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
