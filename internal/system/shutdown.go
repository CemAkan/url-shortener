package system

import (
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"github.com/gofiber/fiber/v2"
	"time"
)

func GracefulShutdown(app *fiber.App) {
	infrastructure.Log.Infof("Starting graceful shutdown...")

	// Fiber shutdown
	if err := app.Shutdown(); err != nil {
		infrastructure.Log.WithError(err).Error("Failed to shutdown Fiber gracefully")
	} else {
		infrastructure.Log.Infof("Fiber shutdown completed")
	}

	// Redis shutdown
	if infrastructure.Redis != nil {
		if err := infrastructure.Redis.Close(); err != nil {
			infrastructure.Log.WithError(err).Error("Failed to close Redis connection")
		} else {
			infrastructure.Log.Infof("Redis connection closed successfully")
		}
	}

	// DB shutdown
	sqlDB, err := infrastructure.DB.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			infrastructure.Log.WithError(err).Error("Failed to close DB pool")
		} else {
			infrastructure.Log.Infof("DB pool closed successfully")
		}
	} else {
		infrastructure.Log.WithError(err).Error("Failed to retrieve DB pool handle")
	}

	infrastructure.Log.Infof("Application shutdown complete. Exiting.")

	//wait to all closings
	infrastructure.Log.Infof("--- Program will close in 10 seconds ---")
	time.Sleep(time.Second * 10)
}
