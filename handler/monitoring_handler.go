package handler

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MonitoringHandlerImpl struct {
	opsProcessed   *prometheus.GaugeVec
	metricsHandler http.Handler
}

func CreateMonitoringHandler() *MonitoringHandlerImpl {
	return &MonitoringHandlerImpl{
		opsProcessed: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "toky_request_time",
			Help: "Request Time records in Seconds",
		}, []string{"URL", "Method"}),
		metricsHandler: promhttp.Handler(),
	}
}

func (m *MonitoringHandlerImpl) MeasureRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.opsProcessed.With(prometheus.Labels{"URL": r.URL.Path, "Method": r.Method}).Set(time.Since(start).Seconds())
	})
}

func (m *MonitoringHandlerImpl) MetricsHandler() http.Handler {
	return m.metricsHandler
}
