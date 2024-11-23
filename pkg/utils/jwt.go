package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTUtil interface {
	GenerateToken(userId uint) (string, error)
	ValidateToken(token string) (uint, error)
}

type jwtUtil struct {
	secretKey []byte
}

func NewJWTUtil(secretKey string) JWTUtil {
	return &jwtUtil{secretKey: []byte(secretKey)}
}

func (j *jwtUtil) GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtUtil) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := uint(claims["user_id"].(float64))
		return userId, nil
	}

	return 0, err
}
