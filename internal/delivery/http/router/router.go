package router

import (
	_ "github.com/CemAkan/url-shortener/docs"
	"github.com/CemAkan/url-shortener/internal/delivery/http/handler"
	middleware2 "github.com/CemAkan/url-shortener/internal/delivery/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, urlHandler *handler.URLHandler, adminHandler *handler.AdminHandler, verificationHandler *handler.VerificationHandler) {
	api := app.Group("/api")

	// Swagger UI Route
	api.Get("/docs/*", fiberSwagger.WrapHandler)

	// implement log middleware
	app.Use(middleware2.RequestLogger())

	// public routes (no need jwt)

	//mail assets
	api.Static("/assets", "app/email/assets")

	//health
	api.Get("/health", handler.Health)

	//redirect
	app.Get("/:code", urlHandler.Redirect)

	//verification
	api.Get("/verify/mail/:token", middleware2.JWTVerification("email_verification"), verificationHandler.VerifyMailAddress)
	api.Get("/verify/password/:token", middleware2.JWTVerification("password_reset_verification"), verificationHandler.ResetPasswordTokenResolve)
	api.Post("/verify/password", middleware2.JWTAuth("password_reset_verification"), verificationHandler.ResetPassword)

	//auth
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// -- protected routes (jwt required) --

	//user
	api.Get("/me", middleware2.JWTAuth("auth"), authHandler.Me)
	api.Get("/password/reset", middleware2.JWTAuth("auth"), authHandler.ResetPassword)

	//url
	api.Post("/shorten", middleware2.JWTAuth("auth"), urlHandler.Shorten)
	api.Get("/my/urls", middleware2.JWTAuth("auth"), urlHandler.ListUserURLs)
	api.Get("/my/urls/:code", middleware2.JWTAuth("auth"), urlHandler.GetSingleURL)
	api.Delete("/my/urls/:code", middleware2.JWTAuth("auth"), urlHandler.DeleteURL)
	api.Patch("/my/urls/:code", middleware2.JWTAuth("auth"), urlHandler.UpdateURL)

	// Admin routes
	adminGroup := api.Group("/admin", middleware2.JWTAuth("auth"), middleware2.AdminOnly())

	adminGroup.Get("/users", adminHandler.ListUsers)
	adminGroup.Delete("/users/:id", adminHandler.RemoveUser)
}
