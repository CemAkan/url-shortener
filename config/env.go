package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println("No local .env file found, continue with system env vars.")
	}
}

// GetEnv gets env value with key parameter and returning fallback value if env record is not exist
func GetEnv(key, fallback string) string {
	value, isExist := os.LookupEnv(key)

	if !isExist {
		return fallback
	}

	return value

}
