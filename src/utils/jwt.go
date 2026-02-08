package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nama     string `json:"nama"`
	RoleID   *uint  `json:"role_id"`
	Status   string `json:"status"`
	jwt.RegisteredClaims
}

// GenerateToken generates JWT token
func GenerateToken(userID uint, username, nama string, roleID *uint, status string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key-change-this-in-production"
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Nama:     nama,
		RoleID:   roleID,
		Status:   status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// VerifyToken verifies JWT token and returns claims
func VerifyToken(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key-change-this-in-production"
	}

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
