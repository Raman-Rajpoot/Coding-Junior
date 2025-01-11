package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims struct holds the JWT claims (user-specific data)
type Claims struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new access token with a 1-hour expiration
func GenerateAccessToken(userName string, email string) (string, error) {
	claims := Claims{
		UserName: userName,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Register",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set in environment variables")
	}

	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken generates a new refresh token with a 7-day expiration
func GenerateRefreshToken(userName string, email string) (string, error) {
	claims := Claims{
		UserName: userName,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Register",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
		},
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set in environment variables")
	}

	return token.SignedString([]byte(secret))
}

// ValidateToken validates a JWT token and extracts the claims
func ValidateToken(tokenString string, tokenType string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set in environment variables")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Optionally check if the token type is correct (access or refresh token)
	if tokenType == "access" && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("access token has expired")
	}

	if tokenType == "refresh" && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token has expired")
	}

	return claims, nil
}
