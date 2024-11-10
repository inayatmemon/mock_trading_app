package auth

import (
	"9DotTechnology/trading_platform/constants/common_constants"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenMetadata struct {
	ID          string `json:"userId"`
	Email       string `json:"email"`
	CountryCode string `json:"countryCode"`
	Mobile      string `json:"mobile"`
}

// Secret key used to sign and verify JWT tokens
var jwtSecret = []byte("secretKey")

// GenerateToken generates a new JWT token with a specified metadata and expiration time
func GenerateToken(md TokenMetadata) (string, error) {
	// Define token claims
	claims := jwt.MapClaims{
		"email":       md.Email,
		"countryCode": md.CountryCode,
		"mobile":      md.Mobile,
		"userId":      md.ID,
		"exp":         time.Now().Add(time.Hour * common_constants.TOKEN_EXPIRY_HOURS).Unix(), // Token expiration (1 hour)
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// DecodeToken parses and validates the token, returning the claims if valid
func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// Check if token is valid and return claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Printf("claims: %v\n", claims)
		return claims, nil
	} else {
		return nil, err
	}
}
