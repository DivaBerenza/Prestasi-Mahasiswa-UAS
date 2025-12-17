package utils

import (
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"UAS/app/model"


	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte(os.Getenv("API_KEY"))

// ----------------- JWT -----------------
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
	// Cek apakah token ada di blacklist
	if IsBlacklisted(tokenStr) {
		return nil, errors.New("token has been logged out")
	}

	// Parse token dengan claims yang sesuai
	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Jangan cek permissions di sini, hanya kembalikan claims
	return claims, nil
}


// ----------------- BLACKLIST -----------------
var blacklist = make(map[string]time.Time)
var mu sync.Mutex

func AddToBlacklist(token string, exp time.Time) {
	mu.Lock()
	defer mu.Unlock()
	blacklist[token] = exp
}

func IsBlacklisted(token string) bool {
	mu.Lock()
	defer mu.Unlock()
	exp, exists := blacklist[token]
	if !exists {
		return false
	}
	if time.Now().After(exp) {
		delete(blacklist, token)
		return false
	}
	return true
}

// ----------------- UTILITY -----------------
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("missing token")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}