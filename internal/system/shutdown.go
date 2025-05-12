package system

import (
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	infrastructure2 "github.com/CemAkan/url-shortener/pkg/infrastructure"
	"github.com/gofiber/fiber/v2"
)

func GracefulShutdown(app *fiber.App) {
	infrastructure2.Log.Info("Starting graceful shutdown...")

	// Fiber shutdown
	if err := app.Shutdown(); err != nil {
		infrastructure2.Log.WithError(err).Error("Failed to shutdown Fiber gracefully")
	} else {
		infrastructure2.Log.Info("Fiber shutdown completed")
	}

	// Redis shutdown
	if infrastructure2.Redis != nil {
		if err := infrastructure2.Redis.Close(); err != nil {
			infrastructure2.Log.WithError(err).Error("Failed to close Redis connection")
		} else {
			infrastructure2.Log.Info("Redis connection closed successfully")
		}
	}

	// DB shutdown
	sqlDB, err := infrastructure.DB.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			infrastructure2.Log.WithError(err).Error("Failed to close DB pool")
		} else {
			infrastructure2.Log.Info("DB pool closed successfully")
		}
	} else {
		infrastructure2.Log.WithError(err).Error("Failed to retrieve DB pool handle")
	}

	infrastructure2.Log.Info("Application shutdown complete. Exiting.")
}
