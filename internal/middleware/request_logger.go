package middleware

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/gofiber/fiber/v2"
	"io"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RequestLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     io.MultiWriter(config.Log.Out), // Logrus' output ( /logs/app.log file + stdout )
	})
}
