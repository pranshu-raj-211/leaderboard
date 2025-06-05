package metrics

import (
	"fmt"
	"leaderboard/src/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GameSubmissions = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "game_submissions_total",
			Help: "Total number of submitted game scores",
		},
	)

	ActiveSSEConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sse_connections",
			Help: "Number of actively connected SSE clients",
		},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	HTTPErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP error responses",
		},
		[]string{"method", "endpoint", "code"},
	)

	SSEMessagesSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sse_messages_sent_total",
			Help: "Total number of SSE messages sent",
		},
	)

	RedisLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redis_latency_seconds",
			Help:    "Latency for Redis operations",
			Buckets: prometheus.DefBuckets,
		},
	)

	LeaderboardUpdateDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "leaderboard_update_duration_seconds",
			Help: "Time taken to compute leaderboard updates",
		},
	)

	DroppedSSEConnections = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dropped_sse_connections_total",
			Help: "Total number of dropped SSE connections",
		},
	)
)

func InitMetrics() {
	metrics := []prometheus.Collector{
		GameSubmissions,
		ActiveSSEConnections,
		HTTPRequestDuration,
		HTTPRequests,
		HTTPErrors,
		SSEMessagesSent,
		RedisLatency,
		LeaderboardUpdateDuration,
		DroppedSSEConnections,
	}
	for _, m := range metrics {
		if err := prometheus.Register(m); err != nil {
			config.Error("Metric registration failed", map[string]any{"Metric": m, "Error": err})
		}
	}
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		HTTPRequestDuration.WithLabelValues(c.Request.Method, c.FullPath(), fmt.Sprintf("%d", status)).Observe(duration)
		HTTPRequests.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
		if status >= 400 {
			HTTPErrors.WithLabelValues(c.Request.Method, c.FullPath(), fmt.Sprintf("%d", status)).Inc()
		}
	}
}
