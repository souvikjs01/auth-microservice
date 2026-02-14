package metric

import (
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	IncHits(status int, method, path string, observeTime time.Duration)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type PrometheusMetric struct {
	HitsTotal    prometheus.Counter
	Hits         *prometheus.CounterVec
	ResponseTime *prometheus.HistogramVec
}

func CreateMetrics(address, name string) (Metrics, error) {
	var metr PrometheusMetric
	metr.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})

	if err := prometheus.Register(metr.HitsTotal); err != nil {
		return nil, err
	}

	metr.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "_hits",
	}, []string{"status", "method", "path"})

	if err := prometheus.Register(metr.Hits); err != nil {
		return nil, err
	}

	metr.ResponseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: name + "_response_time",
	}, []string{"status", "method", "path"})

	if err := prometheus.Register(metr.ResponseTime); err != nil {
		return nil, err
	}

	if err := prometheus.Register(collectors.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	go func() {
		router := echo.New()

		router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		if err := router.Start(address); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	return &metr, nil
}

func (m *PrometheusMetric) IncHits(status int, method, path string, observeTime time.Duration) {
	m.HitsTotal.Inc()
	m.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (m *PrometheusMetric) ObserveResponseTime(status int, method, path string, observeTime float64) {
	m.ResponseTime.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}
