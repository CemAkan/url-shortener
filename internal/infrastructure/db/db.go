package db

import (
	"fmt"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/domain/entity"
	"github.com/CemAkan/url-shortener/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetEnv("POSTGRES_HOST", "pg-shortener"),
		config.GetEnv("POSTGRES_USER", "username"),
		config.GetEnv("POSTGRES_PASSWORD", "password"),
		config.GetEnv("POSTGRES_DB", "urlShortener"),
		config.GetEnv("POSTGRES_PORT", "5432"),
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
