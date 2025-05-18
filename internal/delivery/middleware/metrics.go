package middleware

import (
	"strconv"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/CemAkan/url-shortener/internal/metrics"
)

// SetupPrometheus sets prometheus custom middleware
func SetupPrometheus(app *fiber.App, registry *prometheus.Registry) {
	prom := fiberprometheus.NewWithRegistry(
		registry,
		"url_shortener", // namespace
		"http",          // subsystem
		"service",       // label name
		nil,
	)
	prom.RegisterAt(app, "/metrics")

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		status := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()

		metrics.HTTPRequestTotal.WithLabelValues(status, method).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(status, method).Observe(duration)

		return err
	})
}
