package main

import (
	"context"
	"github.com/CemAkan/url-shortener/config"
	appModule "github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/delivery"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	config.LoadEnv()
	config.InitLogger()
	config.InitDB()
	config.InitRedis()

	app := fiber.New()

	// Dependency injection

	//user
	userRepo := repository.NewUserRepository()
	userService := appModule.NewUserService(userRepo)
	userHandler := delivery.NewAuthHandler(userService)

	//url
	urlRepo := repository.NewURLRepository()
	urlService := appModule.NewURLService(urlRepo)
	urlHandler := delivery.NewURLHandler(urlService)

	// Routes
	delivery.SetupRoutes(app, userHandler, urlHandler)

	// Graceful shutdown signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	port := config.GetEnv("PORT", "3000")

	// start fiber server goroutine
	go func() {
		err := app.Listen(":" + port)
		if err != nil {
			config.Log.Fatalf("Port can not listen: %v", err)
			return
		}
		config.Log.Info("Server starting on port: %s", port)

	}()

	// Wait for Ctrl+C
	<-quit
	config.Log.Info("Shutdown signal received, closing server...")

	// Graceful shutdown timeout context
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// fiber shutdown
	if err := app.Shutdown(); err != nil {
		config.Log.WithError(err).Fatal("Failed to shutdown server gracefully")
	} else {
		config.Log.Info("Server shutdown completed successfully")
	}

	// Redis shutdown
	if config.Redis != nil {
		config.Log.Info("Closing Redis connection...")
		if err := config.Redis.Close(); err != nil {
			config.Log.WithError(err).Error("Failed to close Redis connection")
		} else {
			config.Log.Info("Redis connection closed successfully")
		}
	}

	// DB shutdown
	sqlDB, err := config.DB.DB()
	if err != nil {
		config.Log.WithError(err).Error("Failed to retrieve sql.DB from GORM")
	} else {
		config.Log.Info("Closing DB connection pool...")
		if err := sqlDB.Close(); err != nil {
			config.Log.WithError(err).Error("Failed to close DB connection pool")
		} else {
			config.Log.Info("DB connection pool closed successfully")
		}
	}

	config.Log.Info("Application shutdown complete. Exiting.")
}
