package server

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/email/delivery/rabbitmq"
	emailRepository "github.com/souvikjs01/auth-microservice/internal/email/repository"
	emailUsecase "github.com/souvikjs01/auth-microservice/internal/email/usecase"
	"github.com/souvikjs01/auth-microservice/internal/interceptors"
	metric "github.com/souvikjs01/auth-microservice/internal/metric"
	sessionRepository "github.com/souvikjs01/auth-microservice/internal/session/repository"
	sessionUseCase "github.com/souvikjs01/auth-microservice/internal/session/usecase"
	authServergRPC "github.com/souvikjs01/auth-microservice/internal/user/delivery/grpc/service"
	userRepository "github.com/souvikjs01/auth-microservice/internal/user/repository"
	"github.com/souvikjs01/auth-microservice/internal/user/usecase"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	userProtoService "github.com/souvikjs01/auth-microservice/proto"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// grpc auth server
type Server struct {
	logger   logger.Logger
	cfg      *config.Config
	db       *sqlx.DB
	redis    *redis.Client
	amqpConn *amqp.Connection
}

func NewAuthServer(logger logger.Logger, cfg *config.Config, db *sqlx.DB, redis *redis.Client, amqpConn *amqp.Connection) *Server {
	return &Server{
		logger:   logger,
		cfg:      cfg,
		db:       db,
		redis:    redis,
		amqpConn: amqpConn,
	}
}

func (s *Server) Run() error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.Url, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("metric.CreateMetrics: %v", err)
		return err
	}
	s.logger.Infof("Metrics available url: %v, serviceName: %v", s.cfg.Metrics.Url, s.cfg.Metrics.ServiceName)
	im := interceptors.NewInterceptorManager(s.logger, s.cfg, metrics)
	userPgRepo := userRepository.NewUserRepository(s.db)
	userRedisRepo := userRepository.NewUserRedisRepository(s.redis, s.logger)
	useUC := usecase.NewUserUsecase(s.logger, userPgRepo, userRedisRepo)
	sessionRepo := sessionRepository.NewSessionRepository(s.redis, s.cfg)
	sessionUC := sessionUseCase.NewSessionUseCase(sessionRepo, s.cfg)

	emailRepo := emailRepository.NewEmailRepository()
	emailUC := emailUsecase.NewEmailUsecase(emailRepo, s.logger)
	emailsAmqpConsumer := rabbitmq.NewEmailConsumer(s.amqpConn, s.logger, emailUC)

	go func() {
		err := emailsAmqpConsumer.StartConsumer(
			s.cfg.RabbitMQ.WorkerPoolSize,
			s.cfg.RabbitMQ.Exchange,
			s.cfg.RabbitMQ.Queue,
			s.cfg.RabbitMQ.RoutingKey,
			s.cfg.RabbitMQ.ConsumerTag,
		)

		if err != nil {
			s.logger.Errorf("emailsAmqpConsumer.StartConsumer: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}
	defer lis.Close()

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.Server.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.Server.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.Server.MaxConnectionAge * time.Minute,
		Time:              s.cfg.Server.Time * time.Minute,
	}),
		grpc.ChainUnaryInterceptor(im.Logger, im.Metrics, grpcrecovery.UnaryServerInterceptor(), grpc_prometheus.UnaryServerInterceptor),
	)

	if s.cfg.Server.Mode != "Production" {
		reflection.Register(server)
	}

	authGrpcServer := authServergRPC.NewAuthServerGrpc(s.logger, s.cfg, useUC, sessionUC, metrics)
	userProtoService.RegisterUserServiceServer(server, authGrpcServer)

	grpc_prometheus.Register(server)
	http.Handle("/metrics", promhttp.Handler())
	s.logger.Infof("Server is listening on port: %v", s.cfg.Server.Port)
	if err := server.Serve(lis); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)
	<-quit

	s.logger.Info("Server is shutting down, exited properly")

	return nil
}
