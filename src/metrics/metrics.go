package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    GameSubmissions = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "game_submissions_total",
            Help: "Total number of submitted games",
        },
    )

    ActiveSSEConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_sse_connections",
            Help: "Number of currently connected SSE clients",
        },
    )
)

func InitMetrics() {
    prometheus.MustRegister(GameSubmissions, ActiveSSEConnections)
}