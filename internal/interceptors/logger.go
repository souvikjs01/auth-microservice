package interceptors

import (
	"context"
	"time"

	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorManager struct {
	logger logger.Logger
	cfg    *config.Config
}

func NewInterceptorManager(logger logger.Logger, cfg *config.Config) *InterceptorManager {
	return &InterceptorManager{
		logger: logger,
		cfg:    cfg,
	}
}

func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)

	im.logger.Infof("METHOD: %s, REQUEST: %v, RESPONSE: %v, ERROR: %v, TIME: %s, METADATA: %v", info.FullMethod, req, reply, err, time.Since(start), md)

	return reply, nil
}
