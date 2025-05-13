package config

import (
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

func ResolveToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
