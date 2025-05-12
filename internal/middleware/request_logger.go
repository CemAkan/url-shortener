package middleware

import (
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RequestLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     infrastructure.ServerLogOutput,
	})
}
