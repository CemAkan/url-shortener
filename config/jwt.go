package config

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte(GetEnv("JWT_SECRET", "default_jwt_secret123!@#"))

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)

}
