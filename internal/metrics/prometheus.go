package metrics

import (
	"github.com/CemAkan/url-shortener/internal/delivery/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (

	// DB Health
	DBUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_up",
		Help: "Database connection status (1 = up, 0 = down)",
	})

	// Redis Health
	RedisUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "redis_up",
		Help: "Redis connection status (1 = up, 0 = down)",
	})

	//Mail Health
	MailUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mail_up",
		Help: "Mail service connection status (1 = up, 0 = down)",
	})

	// URL shortening latency
	RedirectLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "shortener_shorten_latency_seconds",
		Help:    "Latency histogram for URL redirecting",
		Buckets: prometheus.DefBuckets,
	})

	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests processed, labeled by status code",
		},
		[]string{"status_code"}, // label listesi burada
	)
)

var alreadyRegistered bool

func RegisterAll() {
	if alreadyRegistered {
		return
	}
	alreadyRegistered = true

	middleware.CustomRegistry.MustRegister(
		DBUp,
		RedisUp,
		MailUp,
		HTTPRequestTotal,
		RedirectLatency,
	)
}
