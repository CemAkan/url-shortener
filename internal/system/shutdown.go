package system

import (
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"github.com/gofiber/fiber/v2"
)

func GracefulShutdown(app *fiber.App) {
	infrastructure.Log.Info("Starting graceful shutdown...")

	// Fiber shutdown
	if err := app.Shutdown(); err != nil {
		infrastructure.Log.WithError(err).Error("Failed to shutdown Fiber gracefully")
	} else {
		infrastructure.Log.Info("Fiber shutdown completed")
	}

	// Redis shutdown
	if infrastructure.Redis != nil {
		if err := infrastructure.Redis.Close(); err != nil {
			infrastructure.Log.WithError(err).Error("Failed to close Redis connection")
		} else {
			infrastructure.Log.Info("Redis connection closed successfully")
		}
	}

	// DB shutdown
	sqlDB, err := infrastructure.DB.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			infrastructure.Log.WithError(err).Error("Failed to close DB pool")
		} else {
			infrastructure.Log.Info("DB pool closed successfully")
		}
	} else {
		infrastructure.Log.WithError(err).Error("Failed to retrieve DB pool handle")
	}

	infrastructure.Log.Info("Application shutdown complete. Exiting.")
}
