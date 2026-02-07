package main

import (
	"fmt"
	"log"
	"os"

	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"github.com/souvikjs01/auth-microservice/pkg/utils"
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
}
