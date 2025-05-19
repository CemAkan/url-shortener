// @title URL Shortener API
// @version 1.0
// @description Enterprise grade URL shortening service.
// @host localhost:3000
// @BasePath /api

package main

import (
	"context"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/delivery/http/handler"
	"github.com/CemAkan/url-shortener/internal/delivery/http/router"
	"github.com/CemAkan/url-shortener/internal/health"
	"github.com/CemAkan/url-shortener/internal/infrastructure/cache"
	"github.com/CemAkan/url-shortener/internal/infrastructure/db"
	"github.com/CemAkan/url-shortener/internal/infrastructure/mail"
	"github.com/CemAkan/url-shortener/internal/jobs"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/CemAkan/url-shortener/internal/seed"
	"github.com/CemAkan/url-shortener/internal/service"
	"github.com/CemAkan/url-shortener/internal/system"
	"github.com/CemAkan/url-shortener/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

func main() {
	config.LoadEnv()
	mail.InitMail()
	logger.InitLogger()
	db.InitDB()
	seed.SeedAdminUser()
	cache.InitRedis()

	appFiber := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		StrictRouting:         false,
		CaseSensitive:         true,
		ProxyHeader:           string(fiber.HeaderXForwardedFor),
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          10 * time.Second,
		TrustedProxies:        strings.Split(config.GetEnv("TRUST_PROXY_IPS", "0.0.0.0/0"), ","),
	})

	// Dependency injection

	//MAIL
	mailService := service.NewMailService()

	//USER

	//auth
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userService, mailService)

	//verification
	verificationHandler := handler.NewVerificationHandler(userService)

	//URL
	urlRepo := repository.NewURLRepository()
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	//ADMIN
	adminHandler := handler.NewAdminHandler(userService, urlService)

	//send handlers
	router.SetupRoutes(appFiber, authHandler, urlHandler, adminHandler, verificationHandler)

	//jobs
	clickFlusher := service.NewClickFlusherService(urlRepo)
	go job.StartClickFlushJob(clickFlusher, 1*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//single handler start
	go system.HandleSignals(cancel)

	//health watchdog start
	go health.StartWatchdog(ctx)

	//server start
	go startServer(appFiber, cancel)

	<-ctx.Done() //chan-block wait

	//graceful shutdown start
	system.GracefulShutdown(appFiber)
}

func startServer(app *fiber.App, cancel context.CancelFunc) {
	port := config.GetEnv("APP_PORT", "3000")
	logger.Log.Infof("Starting Fiber on port: %s", port)

	err := app.Listen(":" + port)

	if err != nil {
		logger.Log.WithError(err).Error("Fiber server failed to start")
		cancel()
	}
}
