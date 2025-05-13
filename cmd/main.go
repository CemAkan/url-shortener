// @title URL Shortener API
// @version 1.0
// @description Enterprise grade URL shortening service.
// @host localhost:3000
// @BasePath /api

package main

import (
	"context"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/delivery"
	"github.com/CemAkan/url-shortener/internal/health"
	job "github.com/CemAkan/url-shortener/internal/jobs"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/CemAkan/url-shortener/internal/system"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"github.com/gofiber/fiber/v2"
	"time"
)

func main() {
	config.LoadEnv()
	infrastructure.InitLogger()
	infrastructure.InitDB()
	infrastructure.InitRedis()

	appFiber := fiber.New()

	// Dependency injection
	userRepo := repository.NewUserRepository()
	userService := app.NewUserService(userRepo)
	userHandler := delivery.NewAuthHandler(userService)

	urlRepo := repository.NewURLRepository()
	urlService := app.NewURLService(urlRepo)
	urlHandler := delivery.NewURLHandler(urlService)

	adminHandler := delivery.NewAdminHandler(userService, urlService)

	delivery.SetupRoutes(appFiber, userHandler, urlHandler, adminHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//jobs
	clickFlusher := app.NewClickFlusherService(urlRepo)
	go job.StartClickFlushJob(clickFlusher, 1*time.Minute)

	go system.HandleSignals(cancel)
	go health.StartWatchdog(ctx)

	go startServer(appFiber, cancel)

	<-ctx.Done()
	defer system.GracefulShutdown(appFiber)
}

func startServer(app *fiber.App, cancel context.CancelFunc) {
	port := config.GetEnv("PORT", "3000")
	infrastructure.Log.Infof("Starting Fiber on port: %s", port)

	err := app.Listen(":" + port)

	if err != nil {
		infrastructure.Log.WithError(err).Error("Fiber server failed to start")
		cancel()
	}
}
