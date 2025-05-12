package infrastructure

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		Log.WithError(err).Fatal("Failed to connect to database")
	}

	Log.Info("Database connection established successfully")
	return db
}
