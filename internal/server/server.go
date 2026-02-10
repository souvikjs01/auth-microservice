package server

import (
	"net"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/interceptors"
	authServergRPC "github.com/souvikjs01/auth-microservice/internal/user/delivery/grpc/server"
	"github.com/souvikjs01/auth-microservice/internal/user/repository"
	"github.com/souvikjs01/auth-microservice/internal/user/usecase"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	userService "github.com/souvikjs01/auth-microservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// grpc auth server
type Server struct {
	logger logger.Logger
	cfg    *config.Config
	db     *sqlx.DB
	redis  *redis.Client
}

func NewAuthServer(logger logger.Logger, cfg *config.Config, db *sqlx.DB, redis *redis.Client) *Server {
	return &Server{
		logger: logger,
		cfg:    cfg,
		db:     db,
		redis:  redis,
	}
}

func (s *Server) Run() error {
	im := interceptors.NewInterceptorManager(s.logger, s.cfg)
	userRepo := repository.NewUserRepository(s.db)
	useUC := usecase.NewUserUsecase(s.logger, userRepo)

	lis, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
		Timeout:           15 * time.Second,
		MaxConnectionAge:  5 * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(grpcrecovery.UnaryServerInterceptor()),
	)

	if s.cfg.Server.Mode != "Production" {
		reflection.Register(server)
	}

	authGrpcServer := authServergRPC.NewAuthServerGrpc(s.logger, s.cfg, useUC)
	userService.RegisterUserServiceServer(server, authGrpcServer)

	s.logger.Infof("Server is listening on port: %v", s.cfg.Server.Port)
	if err := server.Serve(lis); err != nil {
		return err
	}

	return nil
}
