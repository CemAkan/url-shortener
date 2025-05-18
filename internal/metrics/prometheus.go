package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// Health metrics
	DBUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_up", Help: "Database connection status (1 = up, 0 = down)",
	})
	RedisUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "redis_up", Help: "Redis connection status (1 = up, 0 = down)",
	})
	MailUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mail_up", Help: "Mail service connection status (1 = up, 0 = down)",
	})

	// HTTP metrics
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests processed, labeled by status code and method",
		},
		[]string{"status_code", "method"},
	)
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status_code", "method"},
	)
)

// RegisterAll, verilen Registry üzerinde tüm metrikleri Tekrar kayda engelleyecek şekilde register eder.
func RegisterAll(registry *prometheus.Registry) {
	registry.MustRegister(
		DBUp,
		RedisUp,
		MailUp,
		HTTPRequestTotal,
		HTTPRequestDuration,
	)
}
