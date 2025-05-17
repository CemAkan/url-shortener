package middleware

import (
	logger2 "github.com/CemAkan/url-shortener/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	logFileName = "server"
)

func RequestLogger() fiber.Handler {
	file, err := logger2.FileOpener(logFileName)

	if err != nil {
		file = logger2.Log.Out
	}
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     file,
	})
}
