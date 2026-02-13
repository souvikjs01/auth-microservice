package server

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/internal/models"
	"github.com/souvikjs01/auth-microservice/pkg/grpc_errors"
	"github.com/souvikjs01/auth-microservice/pkg/utils"
	userService "github.com/souvikjs01/auth-microservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

// Login User
func (u *userServer) Login(c context.Context, r *userService.LoginRequest) (*userService.LoginResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "user.Login")
	defer span.Finish()

	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for k, v := range incomingContext {
			log.Printf("key: %s, value: %v", k, v)
		}
	}

	email := r.GetEmail()

	if !utils.ValidateEmail(email) {
		u.logger.Errorf("invalid email: %s", email)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(errors.New("ErrInvalid email address")), "invalid email: %s", email)
	}

	user, err := u.userUC.Login(ctx, email, r.GetPassword())
	if err != nil {
		u.logger.Errorf("userUC.Login %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "userUC.Login: %v", err)
	}

	session, err := u.sessionUC.CreateSession(ctx, &models.Session{
		UserID: user.UserID.String(),
	}, u.cfg.Session.Expire)

	if err != nil {
		u.logger.Errorf("sessionUC.CreateSession: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "sessionUC.CreateSession: %v", err)
	}

	return &userService.LoginResponse{
		User:      u.userModelToProto(user),
		SessionId: session,
	}, nil
}

// get session id from context and find user by id and return user
func (u *userServer) GetMe(ctx context.Context, r *userService.GetMeRequest) (*userService.GetMeResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.GetMe")
	defer span.Finish()

	sessionID, err := u.getSessionIDFromContext(ctx)
	if err != nil {
		u.logger.Errorf("getSessionIDFromContext: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "no context data received")
	}

	session, err := u.sessionUC.GetSessionID(ctx, sessionID)
	if err != nil {
		u.logger.Errorf("sessionUC.GetSessionID: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "sessionUC.GetSessionID: %v", err)
	}

	user, err := u.userUC.FindByID(ctx, uuid.MustParse(session.UserID))
	if err != nil {
		u.logger.Errorf("userUC.FindByID: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "userUC.FindByID: %v", err)
	}
	return &userService.GetMeResponse{
		User: u.userModelToProto(user),
	}, nil
}

// Logout user and delete current session
func (u *userServer) Logout(ctx context.Context, r *userService.LogoutRequest) (*userService.LogoutResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Logout")
	defer span.Finish()

	sessionID, err := u.getSessionIDFromContext(ctx)
	if err != nil {
		u.logger.Errorf("getSessionIDFromContext: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "no context data received")
	}

	if err := u.sessionUC.DeleteByID(ctx, sessionID); err != nil {
		u.logger.Errorf("sessionUC.DeleteByID: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "sessionUC.DeleteByID: %v", err)
	}

	return &userService.LogoutResponse{}, nil
}

// helper methods
func (u *userServer) getSessionIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "No ctx data")
	}

	sessionID := md.Get("session_id")
	if sessionID[0] == "" {
		return "", status.Error(codes.Unauthenticated, "No session_id in ctx metadata")
	}

	return sessionID[0], nil
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
