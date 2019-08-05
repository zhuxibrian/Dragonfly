package server

import (
	"net/http"

	"github.com/dragonflyoss/Dragonfly/pkg/util"
	"github.com/dragonflyoss/Dragonfly/supernode/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// metrics defines three prometheus metrics for monitoring http handler status
type metrics struct {
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
}

func newMetrics() *metrics {
	return &metrics{
		requestCounter: util.NewCounter(config.SubsystemSupernode, "http_requests_total",
			"Counter of HTTP requests.", []string{"code", "handler", "method"},
		),
		requestDuration: util.NewHistogram(config.SubsystemSupernode, "http_request_duration_seconds",
			"Histogram of latencies for HTTP requests.", []string{"code", "handler", "method"},
			[]float64{.1, .2, .4, 1, 3, 8, 20, 60, 120},
		),
		requestSize: util.NewHistogram(config.SubsystemSupernode, "http_request_size_bytes",
			"Histogram of request size for HTTP requests.", []string{"code", "handler", "method"},
			prometheus.ExponentialBuckets(100, 10, 8),
		),
		responseSize: util.NewHistogram(config.SubsystemSupernode, "http_response_size_bytes",
			"Histogram of response size for HTTP requests.", []string{"code", "handler", "method"},
			prometheus.ExponentialBuckets(100, 10, 8),
		),
	}
}

// instrumentHandler will update metrics for every http request
func (m *metrics) instrumentHandler(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
	return promhttp.InstrumentHandlerDuration(
		m.requestDuration.MustCurryWith(prometheus.Labels{"handler": handlerName}),
		promhttp.InstrumentHandlerCounter(
			m.requestCounter.MustCurryWith(prometheus.Labels{"handler": handlerName}),
			promhttp.InstrumentHandlerRequestSize(
				m.requestSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
				promhttp.InstrumentHandlerResponseSize(
					m.responseSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
					handler,
				),
			),
		),
	)
}
