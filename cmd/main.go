package main

import (
	"context"
	"github.com/CemAkan/url-shortener/config"
	appModule "github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/delivery"
	"github.com/CemAkan/url-shortener/internal/health"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Init phase
	config.LoadEnv()
	config.InitLogger()
	config.InitDB()
	config.InitRedis()

	appFiber := fiber.New()

	// Dependency injection
	userRepo := repository.NewUserRepository()
	userService := appModule.NewUserService(userRepo)
	userHandler := delivery.NewAuthHandler(userService)

	urlRepo := repository.NewURLRepository()
	urlService := appModule.NewURLService(urlRepo)
	urlHandler := delivery.NewURLHandler(urlService)

	delivery.SetupRoutes(appFiber, userHandler, urlHandler)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling (SIGINT, SIGTERM)
	go handleSignals(cancel)

	// Watchdog health monitoring
	go health.StartWatchdog(ctx, cancel)

	// Fiber server
	go func() {
		port := config.GetEnv("PORT", "3000")
		config.Log.Infof("Starting Fiber on port: %s", port)
		if err := appFiber.Listen(":" + port); err != nil {
			config.Log.WithError(err).Error("Fiber server failed to start")
			cancel()
		}
	}()

	// Wait for context cancel (signal or watchdog triggers this)
	<-ctx.Done()
	config.Log.Info("Shutdown initiated...")

	gracefulShutdown(appFiber)
}

func handleSignals(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	sig := <-sigs
	config.Log.Infof("Received signal: %s", sig)
	cancel()
}

func gracefulShutdown(app *fiber.App) {
	// Fiber shutdown
	if err := app.Shutdown(); err != nil {
		config.Log.WithError(err).Error("Failed to shutdown Fiber gracefully")
	} else {
		config.Log.Info("Fiber shutdown completed")
	}

	// Redis shutdown
	if config.Redis != nil {
		if err := config.Redis.Close(); err != nil {
			config.Log.WithError(err).Error("Failed to close Redis connection")
		} else {
			config.Log.Info("Redis connection closed")
		}
	}

	// DB shutdown
	sqlDB, err := config.DB.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			config.Log.WithError(err).Error("Failed to close DB pool")
		} else {
			config.Log.Info("DB pool closed successfully")
		}
	} else {
		config.Log.WithError(err).Error("Failed to retrieve DB pool")
	}

	config.Log.Info("Application shutdown complete. Exiting.")
}
