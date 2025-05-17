package system

import (
	"github.com/CemAkan/url-shortener/pkg/infrastructure/cache"
	"github.com/CemAkan/url-shortener/pkg/infrastructure/db"
	"github.com/CemAkan/url-shortener/pkg/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
	"time"
)

func GracefulShutdown(app *fiber.App) {
	logger.Log.Infof("Starting graceful shutdown...")

	// Fiber shutdown
	if err := app.Shutdown(); err != nil {
		logger.Log.WithError(err).Error("Failed to shutdown Fiber gracefully")
	} else {
		logger.Log.Infof("Fiber shutdown completed")
	}

	// Redis shutdown
	if cache.Redis != nil {
		if err := cache.Redis.Close(); err != nil {
			logger.Log.WithError(err).Error("Failed to close Redis connection")
		} else {
			logger.Log.Infof("Redis connection closed successfully")
		}
	}

	// DB shutdown
	sqlDB, err := db.DB.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			logger.Log.WithError(err).Error("Failed to close DB pool")
		} else {
			logger.Log.Infof("DB pool closed successfully")
		}
	} else {
		logger.Log.WithError(err).Error("Failed to retrieve DB pool handle")
	}

	logger.Log.Infof("Application shutdown complete. Exiting.")

	//wait to all closings
	logger.Log.Infof("--- Program will close in 10 seconds ---")
	time.Sleep(time.Second * 10)
}
