package config

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte(GetEnv("JWT_SECRET", "default_jwt_secret123!@#"))

// GenerateJWT creates signed token with user ID and 24h expiration
func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)

}

// GenerateJWT parses and verifies signed tokens
func ValidateJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// signature method checking
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}
