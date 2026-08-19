package utils

import (
	"errors"
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
	AppName  string `json:"app_name"` // Application name: "PINTU" or "SIEKSA"
	jwt.RegisteredClaims
}

// GenerateToken generates JWT token with app-specific secret
func GenerateToken(userID uint, username, nama string, roleID *uint, status, appName string) (string, error) {
	// Select secret key based on application
	var secretKey string
	switch appName {
	case "PINTU":
		secretKey = os.Getenv("JWT_SECRET_PINTU")
		if secretKey == "" {
			secretKey = os.Getenv("JWT_SECRET") // Fallback to old secret for backward compatibility
			if secretKey == "" {
				secretKey = "your-secret-key-change-this-in-production"
			}
		}
	case "SIEKSA":
		secretKey = os.Getenv("JWT_SECRET_SIEKSA")
		if secretKey == "" {
			return "", errors.New("JWT secret not configured for SIEKSA")
		}
	default:
		return "", errors.New("invalid app name: " + appName)
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Nama:     nama,
		RoleID:   roleID,
		Status:   status,
		AppName:  appName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// VerifyToken verifies JWT token with app-specific secret and returns claims
func VerifyToken(tokenString, appName string) (*JWTClaims, error) {
	// Select secret key based on application
	var secretKey string
	switch appName {
	case "PINTU":
		secretKey = os.Getenv("JWT_SECRET_PINTU")
		if secretKey == "" {
			secretKey = os.Getenv("JWT_SECRET") // Fallback to old secret for backward compatibility
			if secretKey == "" {
				secretKey = "your-secret-key-change-this-in-production"
			}
		}
	case "SIEKSA":
		secretKey = os.Getenv("JWT_SECRET_SIEKSA")
		if secretKey == "" {
			return nil, errors.New("JWT secret not configured for SIEKSA")
		}
	default:
		return nil, errors.New("invalid app name: " + appName)
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

	// Validate app_name in token matches the expected app
	if claims.AppName != appName {
		return nil, errors.New("token not valid for this application")
	}

	return claims, nil
}
