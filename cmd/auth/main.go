package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/souvikjs01/auth-microservice/config"
	authServergRPC "github.com/souvikjs01/auth-microservice/internal/auth/delivery/grpc/server"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"github.com/souvikjs01/auth-microservice/pkg/postgres"
	"github.com/souvikjs01/auth-microservice/pkg/redis"
	"github.com/souvikjs01/auth-microservice/pkg/utils"
	userService "github.com/souvikjs01/auth-microservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("hello there")

	configPath := utils.GetConfigPath(os.Getenv("config"))
	fmt.Println(configPath)
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Load config error: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("Parseconfing: %v", err)
	}

	appLogger := logger.NewLogger(cfg)
	appLogger.InitLogger()

	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.Ssl)
	appLogger.Infof("Success parsed config: %v", cfg.Server.AppVersion)

	// postgres conn
	pgsqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Info("Postgresql connected successfully")
	}
	defer pgsqlDB.Close()

	// redis conn
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis is connected")

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		appLogger.Fatal(err)
	}

	server := grpc.NewServer(grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Timeout:           15 * time.Second,
			MaxConnectionAge:  5 * time.Minute,
		},
	))

	if cfg.Server.Mode != "Production" {
		reflection.Register(server)
	}

	authGrpcServer := authServergRPC.NewAuthServerGrpc(appLogger, cfg)
	userService.RegisterUserServiceServer(server, authGrpcServer)

	appLogger.Infof("Server is listening on port: %v", cfg.Server.Port)
	appLogger.Fatal(server.Serve(lis))
}
