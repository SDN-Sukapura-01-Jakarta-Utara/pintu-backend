package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Request counters
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Request duration
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Request size
	requestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_size_bytes",
			Help: "HTTP request size in bytes",
		},
		[]string{"method", "endpoint"},
	)

	// Response size
	responseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_size_bytes",
			Help: "HTTP response size in bytes",
		},
		[]string{"method", "endpoint"},
	)

	// Errors counter
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "endpoint", "status"},
	)
)

// PrometheusMiddleware returns a Gin middleware for Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		method := c.Request.Method
		endpoint := c.Request.URL.Path

		// Get request size
		reqSize := c.Request.ContentLength
		if reqSize > 0 {
			requestSize.WithLabelValues(method, endpoint).Observe(float64(reqSize))
		}

		c.Next()

		// Calculate metrics
		duration := time.Since(startTime).Seconds()
		statusCode := c.Writer.Status()

		// Record metrics
		requestsTotal.WithLabelValues(method, endpoint, string(rune(statusCode))).Inc()
		requestDuration.WithLabelValues(method, endpoint).Observe(duration)
		responseSize.WithLabelValues(method, endpoint).Observe(float64(c.Writer.Size()))

		// Count errors (5xx and 4xx)
		if statusCode >= 400 {
			errorsTotal.WithLabelValues(method, endpoint, string(rune(statusCode))).Inc()
		}
	}
}
