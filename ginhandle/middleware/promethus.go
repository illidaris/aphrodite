package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMiddleware interface {
	// WrapHandler wraps the given HTTP handler for instrumentation.
	WrapHandler(handlerName string, handler http.Handler) http.HandlerFunc
	WrapGinHandler(handlerName string, handler http.Handler) gin.HandlerFunc
}

type prometheusMiddleware struct {
	buckets  []float64
	registry prometheus.Registerer
}

func (m *prometheusMiddleware) WrapGinHandler(handlerName string, handler http.Handler) gin.HandlerFunc {
	srv := m.WrapHandler(handlerName, handler)
	return gin.WrapH(srv)
}

// WrapHandler wraps the given HTTP handler for instrumentation:
// It registers four metric collectors (if not already done) and reports HTTP
// metrics to the (newly or already) registered collectors.
// Each has a constant label named "handler" with the provided handlerName as
// value.
func (m *prometheusMiddleware) WrapHandler(handlerName string, handler http.Handler) http.HandlerFunc {
	reg := prometheus.WrapRegistererWith(prometheus.Labels{"handler": handlerName}, m.registry)

	requestsTotal := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Tracks the number of HTTP requests.",
		}, []string{"method", "code"},
	)
	requestDuration := promauto.With(reg).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Tracks the latencies for HTTP requests.",
			Buckets: m.buckets,
		},
		[]string{"method", "code"},
	)
	requestSize := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "Tracks the size of HTTP requests.",
		},
		[]string{"method", "code"},
	)
	responseSize := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "Tracks the size of HTTP responses.",
		},
		[]string{"method", "code"},
	)

	// Wraps the provided http.Handler to observe the request result with the provided metrics.
	base := promhttp.InstrumentHandlerCounter(
		requestsTotal,
		promhttp.InstrumentHandlerDuration(
			requestDuration,
			promhttp.InstrumentHandlerRequestSize(
				requestSize,
				promhttp.InstrumentHandlerResponseSize(
					responseSize,
					handler,
				),
			),
		),
	)

	return base.ServeHTTP
}

// New returns a Middleware interface.
func NewPrometheusMiddleware(registry prometheus.Registerer, buckets []float64) PrometheusMiddleware {
	if buckets == nil {
		buckets = prometheus.ExponentialBuckets(0.1, 1.5, 5)
	}

	return &prometheusMiddleware{
		buckets:  buckets,
		registry: registry,
	}
}

func PrometheusHandle() gin.HandlerFunc {
	registry := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return NewPrometheusMiddleware(registry, nil).
		WrapGinHandler("/metrics", promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{}),
		)
}
