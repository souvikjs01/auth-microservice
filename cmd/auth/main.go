package main

import (
	"fmt"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/server"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"github.com/souvikjs01/auth-microservice/pkg/postgres"
	"github.com/souvikjs01/auth-microservice/pkg/redis"
	"github.com/souvikjs01/auth-microservice/pkg/utils"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
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
	} else {
		appLogger.Info("Postgresql connected successfully")
	}
	defer pgsqlDB.Close()

	// redis conn
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis is connected")

	// jaeger
	jaegerCfgInstance := jaegerCfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, // sample ALL traces (dev only)
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:          cfg.Jaeger.LogSpans,
			CollectorEndpoint: "http://localhost:14268/api/traces",
			// LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegerCfg.Logger(jaeger.StdLogger),
		jaegerCfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		log.Fatal("can not create tracer", err)
	}

	appLogger.Info("Jaeger is connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	authServer := server.NewAuthServer(appLogger, cfg)
	appLogger.Fatal(authServer.Run())
}
