package interceptors

import (
	"context"
	"net/http"
	"time"

	"github.com/souvikjs01/auth-microservice/config"
	"github.com/souvikjs01/auth-microservice/internal/metric"
	"github.com/souvikjs01/auth-microservice/pkg/grpc_errors"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorManager struct {
	logger logger.Logger
	cfg    *config.Config
	mtr    metric.Metrics
}

func NewInterceptorManager(logger logger.Logger, cfg *config.Config, mtr metric.Metrics) *InterceptorManager {
	return &InterceptorManager{
		logger: logger,
		cfg:    cfg,
		mtr:    mtr,
	}
}

func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)

	im.logger.Infof("METHOD: %s, REQUEST: %v, RESPONSE: %v, ERROR: %v, TIME: %s, METADATA: %v", info.FullMethod, req, reply, err, time.Since(start), md)

	return reply, nil
}

func (im *InterceptorManager) Metrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	status := http.StatusOK
	if err != nil {
		status = grpc_errors.MapGRPCCodeToHTTPStatus(grpc_errors.ParseGRPCErrStatusCode(err))
	}
	im.mtr.ObserveResponseTime(status, info.FullMethod, info.FullMethod, time.Since(start).Seconds())
	im.mtr.IncHits(status, info.FullMethod, info.FullMethod, time.Since(start))
	return resp, err
}
