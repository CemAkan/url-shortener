package middleware

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var CustomRegistry = prometheus.NewRegistry()

func PrometheusMiddleware(app *fiber.App) {
	prometheusMiddleware := fiberprometheus.NewWithRegistry(
		CustomRegistry,
		"url-shortener", // namespace
		"http",          // subsystem
		"service",       // service label name
		map[string]string{
			"env": "dev",
		},
	)

	prometheusMiddleware.RegisterAt(app, "/metrics")
	app.Use(prometheusMiddleware.Middleware)
}
