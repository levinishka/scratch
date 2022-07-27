package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HttpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests.",
	},
	[]string{"path"},
)

var HttpRequestsDurationSeconds = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_requests_duration_seconds",
		Help: "Duration of HTTP requests.",
	},
	[]string{"path"},
)

var HttpResponseStatusCodesTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_response_status_codes_total",
		Help: "Number of response status codes for path",
	},
	[]string{"path", "code"},
)
