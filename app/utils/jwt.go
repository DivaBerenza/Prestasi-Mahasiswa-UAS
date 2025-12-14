package utils

import (
	"errors"
	"os"
	"time"

	"UAS/app/model"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte(os.Getenv("API_KEY")) // pastikan API_KEY di .env

func GenerateJWT(userID, roleID string, perms []string) (string, error) {
	claims := model.JWTClaims{
		UserID:      userID,
		RoleID:      roleID,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func ValidateJWT(tokenStr string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
