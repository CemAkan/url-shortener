package middleware

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

// PrometheusMiddleware registers /metrics endpoint and middleware for metrics tracking
func PrometheusMiddleware(app *fiber.App) {
	prometheus := fiberprometheus.New("url-shortener")
	prometheus.RegisterAt(app, "/api/metrics") // exposes at /metrics
	app.Use(prometheus.Middleware)             // auto-collects metrics
}
