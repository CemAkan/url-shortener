package seed

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/domain/entity"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/CemAkan/url-shortener/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdminUser() {
	userRepo := repository.NewUserRepository()
	adminMail := config.GetEnv("ADMIN_EMAIL", "")
	adminPass := config.GetEnv("ADMIN_PASSWORD", "")

	if adminPass == "" || adminMail == "" {
		logger.Log.Infof("Env not set, skipping admin seeding.")
		return
	}

	if isExist, _ := userRepo.FindByEmail(adminMail); isExist != nil {
		logger.Log.Infof("Admin user already exists. Skipping seeding.")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Fatalf("Failed to hash password: %v", err)
	}

	admin := entity.User{
		Name:            "initial",
		Surname:         "admin",
		Email:           adminMail,
		Password:        string(hashedPassword),
		IsAdmin:         true,
		IsMailConfirmed: true,
	}

	if err := userRepo.Create(&admin); err != nil {
		logger.Log.Fatalf("Failed to create admin user: %v", err)
	}

	logger.Log.Infof("Admin user seeded successfully.")
}
