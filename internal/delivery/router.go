package delivery

import (
	_ "github.com/CemAkan/url-shortener/docs"
	"github.com/CemAkan/url-shortener/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App, authHandler *AuthHandler, urlHandler *URLHandler, adminHandler *AdminHandler, verificationHandler *VerificationHandler) {
	api := app.Group("/api")

	// Swagger UI Route
	api.Get("/docs/*", fiberSwagger.WrapHandler)

	// implement log middleware
	app.Use(middleware.RequestLogger())

	// public routes (no need jwt)

	//health
	api.Get("/health", Health)

	//redirect
	app.Get("/:code", urlHandler.Redirect)

	//verification
	app.Get("/verify/mail/:token", middleware.JWTVerification("email_verification"), verificationHandler.VerifyMailAddress)
	app.Post("/verify/password/:token", middleware.JWTVerification("password_reset_verification"), verificationHandler.ResetPassword)

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
