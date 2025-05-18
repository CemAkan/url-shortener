package middleware

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"

	"github.com/CemAkan/url-shortener/internal/metrics"
	"github.com/gofiber/fiber/v2"
)

var CustomRegistry = prometheus.NewRegistry()

func MetricsMiddleware(app *fiber.App) fiber.Handler {
	prometheusMiddleware := fiberprometheus.NewWithRegistry(CustomRegistry, "url_shortener", "http", "service", nil)
	prometheusMiddleware.RegisterAt(app, "/metrics")

	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()

		metrics.HTTPRequestTotal.WithLabelValues(status, method).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(status, method).Observe(duration)

		return err
	}
}
