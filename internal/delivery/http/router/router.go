package router

import (
	"github.com/CemAkan/url-shortener/config"
	_ "github.com/CemAkan/url-shortener/docs"
	"github.com/CemAkan/url-shortener/internal/delivery/http/handler"
	"github.com/CemAkan/url-shortener/internal/delivery/middleware"
	"github.com/CemAkan/url-shortener/internal/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, urlHandler *handler.URLHandler, adminHandler *handler.AdminHandler, verificationHandler *handler.VerificationHandler) {
	// implement metrics middleware
	registry := prometheus.NewRegistry()

	metrics.RegisterAll(registry)

	middleware.SetupPrometheus(app, registry)

	// metric ip protection check
	if config.GetEnv("METRICS_PROTECT", "true") == "true" {
		app.Use("/metrics", middleware.IPWhitelistMiddleware())
	}

	// Swagger UI Route
	app.Get("api/docs/*", fiberSwagger.WrapHandler)

	// swagger ip protection check
	if config.GetEnv("SWAGGER_PROTECT", "true") == "true" {
		app.Use("api/docs/*", middleware.IPWhitelistMiddleware())
	}

	// implement log middleware
	app.Use(middleware.RequestLogger())

	// -- public routes (no need jwt) --

	//redirect
	app.Get("/:code", urlHandler.Redirect)

	api := app.Group("/api")

	//mail assets
	api.Static("/assets", "./email/assets")

	//health
	api.Get("/health", handler.Health)

	//verification
	api.Get("/verify/mail/:token", middleware.JWTVerification("email_verification"), verificationHandler.VerifyMailAddress)
	api.Get("/verify/password/:token", middleware.JWTVerification("password_reset_verification"), verificationHandler.ResetPasswordTokenResolve)
	api.Post("/verify/password", middleware.JWTAuth("password_reset_verification"), verificationHandler.ResetPassword)

	//auth
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// -- protected routes (jwt required) --

	//user
	api.Get("/me", middleware.JWTAuth("auth"), authHandler.Me)
	api.Get("/password/reset", middleware.JWTAuth("auth"), authHandler.ResetPassword)

	//url
	api.Post("/shorten", middleware.JWTAuth("auth"), urlHandler.Shorten)
	api.Get("/my/urls", middleware.JWTAuth("auth"), urlHandler.ListUserURLs)
	api.Get("/my/urls/:code", middleware.JWTAuth("auth"), urlHandler.GetSingleURL)
	api.Delete("/my/urls/:code", middleware.JWTAuth("auth"), urlHandler.DeleteURL)
	api.Patch("/my/urls/:code", middleware.JWTAuth("auth"), urlHandler.UpdateURL)

	// Admin routes
	adminGroup := api.Group("/admin", middleware.JWTAuth("auth"), middleware.AdminOnly())

	adminGroup.Get("/users", adminHandler.ListUsers)
	adminGroup.Delete("/users/:id", adminHandler.RemoveUser)
}
