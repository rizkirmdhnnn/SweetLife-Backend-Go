package helper

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
)

// JWTSIGNKEY is the key used to sign JWT tokens
var JWTSIGNKEY = []byte(config.ENV.JWTSIGNKEY)

// JWTCustomClaims is a custom JWT claims struct
type JWTCustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

// Helper function to generate claims
func generateJWTClaims(user *models.User, expirationDuration time.Duration) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ID:        user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
}

// GenerateToken generates a new JWT token (both access and refresh)
func generateToken(user *models.User, expirationDuration time.Duration) (string, error) {
	claims := JWTCustomClaims{
		ID:               user.ID,
		RegisteredClaims: generateJWTClaims(user, expirationDuration),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(JWTSIGNKEY)
	return ss, err
}

// GenerateTokenAccessToken generates a new access token with a 30-minute expiration
func GenerateTokenAccessToken(user *models.User) (string, error) {
	return generateToken(user, 30*time.Minute)
}

// GenerateTokenRefresh generates a new refresh token with a 7-day expiration
func GenerateTokenRefresh(user *models.User) (string, error) {
	return generateToken(user, 7*24*time.Hour)
}

// ValidateToken validates a JWT token
func validateToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSIGNKEY, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateTokenAccessToken validates an access token
func ValidateTokenAccessToken(tokenString string) (*JWTCustomClaims, error) {
	return validateToken(tokenString)
}

// ValidateTokenRefresh validates a refresh token
func ValidateTokenRefresh(tokenString string) (*JWTCustomClaims, error) {
	return validateToken(tokenString)
}

// GetClaimsFromToken extracts the claims from a JWT token
func GetClaimsFromToken(tokenString string) (*JWTCustomClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSIGNKEY, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTCustomClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
