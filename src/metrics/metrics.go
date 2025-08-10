package metrics

import (
	"fmt"
	"leaderboard/src/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var redisLatencyBuckets = []float64{0.00001, 0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1}

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
			Buckets: redisLatencyBuckets,
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
			Buckets: redisLatencyBuckets,
		},
	)

	LeaderboardUpdateDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "leaderboard_update_duration_seconds",
			Help:    "Time taken to compute leaderboard updates",
			Buckets: redisLatencyBuckets,
		},
	)

	DroppedSSEConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dropped_sse_connections_total",
			Help: "Total number of dropped SSE connections",
		},
	)

	RedisOperationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_operation_errors_total",
			Help: "Redis operation errors by type",
		},
		[]string{"operation"},
	)

	ConcurrentClients = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "concurrent_clients_total",
			Help: "Number of concurrent clients",
		},
	)

	RedisPayloadSize = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redis_payload_size_bytes",
			Help:    "Size of Redis operation payloads",
			Buckets: []float64{128, 512, 1024, 2048, 4096, 8192},
		},
	)

	RedisOperationLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redis_operation_latency_seconds",
			Help:    "Latency by Redis operation type",
			Buckets: redisLatencyBuckets,
		},
	)

	JSONMarshalDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "json_marshal_duration_seconds",
			Help:    "Time taken for JSON marshaling operations",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01},
		},
	)

	JSONErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "json_operation_errors_total",
			Help: "JSON operation errors by type",
		},
		[]string{"operation"},
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
		RedisOperationErrors,
		ConcurrentClients,
		RedisPayloadSize,
		RedisOperationLatency,
		JSONErrors,
		JSONMarshalDuration,
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
