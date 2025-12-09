package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type MetricsService struct {
}

var httpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP requests counter",
	},
	[]string{"method", "endpoint", "status"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "endpoint", "status"},
)

func NewMetricsService() *MetricsService {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)

	prometheus.MustRegister(collectors.NewBuildInfoCollector())

	return &MetricsService{}
}

func (m MetricsService) IncrementCountHttpRequests(method string, endpoint string, status string) {
	httpRequests.WithLabelValues(method, endpoint, status).Inc()
}

func (m MetricsService) IncrementDurationHttpRequests(method string, endpoint string, status string, duration float64) {
	httpDuration.WithLabelValues(method, endpoint, status).Observe(duration)
}
