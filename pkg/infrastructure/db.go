package infrastructure

import (
	"fmt"
	"github.com/CemAkan/url-shortener/internal/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		Log.Fatalf("db connection error %v", err.Error())
	}

	//auto migration
	err = database.AutoMigrate(&model.URL{}, &model.User{})

	if err != nil {
		Log.Fatalf("db connection error %v", err.Error())
	}

	// set global db var
	DB = database
}
