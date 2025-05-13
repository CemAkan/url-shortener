package infrastructure

import (
	"fmt"
	"github.com/CemAkan/url-shortener/internal/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	//log file
	dbLogFile, err := FileOpener("database")
	if err != nil {
		dbLogFile = mainLogFile
	}

	//logger configuration
	newLogger := logger.New(
		log.New(dbLogFile, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	//db open
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil || database == nil {
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
