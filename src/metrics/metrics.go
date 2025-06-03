package metrics

import (
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

    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Histogram of request durations.",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)

func InitMetrics() {
    prometheus.MustRegister(GameSubmissions, ActiveSSEConnections, requestDuration)
}