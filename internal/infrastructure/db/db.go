package db

import (
	"fmt"
	"github.com/CemAkan/url-shortener/internal/domain/entity"
	"github.com/CemAkan/url-shortener/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_NAME"), os.Getenv("POSTGRES_PORT"),
	)

	//log file
	dbLogFile, err := logger.FileOpener("database")
	if err != nil {
		dbLogFile = logger.MainLogFile
	}

	//logger configuration
	newLogger := gormLogger.New(
		log.New(dbLogFile, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold: time.Second,
			LogLevel:      gormLogger.Warn,
			Colorful:      false,
		},
	)

	//db open
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil || database == nil {
		logger.Log.Fatalf("db connection error %v", err)
	}

	//auto migration
	if err := database.AutoMigrate(&entity.URL{}, &entity.User{}); err != nil {
		logger.Log.Fatalf("db connection error %v", err.Error())
	}

	// set global db var
	DB = database

}
