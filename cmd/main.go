package main

import (
	"github.com/CemAkan/url-shortener/config"
	appModule "github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/delivery"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {

	config.LoadEnv()
	config.InitLogger()
	config.InitDB()
	config.InitRedis()

	app := fiber.New()

	// Dependency injection
	userRepo := repository.NewUserRepository()
	userService := appModule.NewUserService(userRepo)
	userHandler := delivery.NewAuthHandler(userService)

	// Routes
	delivery.SetupRoutes(app, userHandler)

	port := config.GetEnv("PORT", "3000")

	err := app.Listen(":" + port)
	if err != nil {
		config.Log.Fatalf("Port can not listen: %v", err)
		return
	}
}
