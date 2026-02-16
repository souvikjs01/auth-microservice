package main

import (
	"fmt"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/server"
	"github.com/souvikjs01/auth-microservice/pkg/jaeger"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"github.com/souvikjs01/auth-microservice/pkg/postgres"
	"github.com/souvikjs01/auth-microservice/pkg/rabbitmq"
	"github.com/souvikjs01/auth-microservice/pkg/redis"
	"github.com/souvikjs01/auth-microservice/pkg/utils"
)

func main() {
	log.Println("Start the user service")

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
	}

	defer pgsqlDB.Close()

	// redis conn
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis is connected")

	// jaeger
	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatalf("cannot create tracer: %s", err)
	}

	appLogger.Info("Jaeger is connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		appLogger.Fatalf("cannot create rabbitmq connection: %v", err)
	}
	defer amqpConn.Close()

	authServer := server.NewAuthServer(appLogger, cfg, pgsqlDB, redisClient, amqpConn)
	appLogger.Fatal(authServer.Run())
}
