package config

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID uint, expiresIn time.Duration, purpose string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiresIn).Unix(),
		"type":    purpose,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr, purpose string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != purpose {
		return 0, errors.New("invalid token type")
	}

	userID := uint(claims["user_id"].(float64))
	return userID, nil
}
