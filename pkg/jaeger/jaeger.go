package jaeger

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/souvikjs01/auth-microservice/config"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// init jaeger
func InitJaeger(cfg *config.Config) (opentracing.Tracer, io.Closer, error) {
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
	return jaegerCfgInstance.NewTracer(
		jaegerCfg.Logger(jaeger.StdLogger),
		jaegerCfg.Metrics(metrics.NullFactory),
	)
}
