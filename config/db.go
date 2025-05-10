package config

import (
	"fmt"
	"github.com/CemAkan/url-shortener/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	// complete dsn url with getting values from env
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		GetEnv("DB_HOST", ""), GetEnv("DB_USER", ""), GetEnv("DB_PASSWORD", ""), GetEnv("DB_NAME", ""), GetEnv("DB_PORT", ""), GetEnv("DB_SSLMODE", ""))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		Log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	Log.Info("Database connection established successfully")

	//Auto Migration
	err = DB.AutoMigrate(&domain.User{})
	if err != nil {
		Log.Fatal("❌ DB AutoMigrate failed: ", err)
	}
}
