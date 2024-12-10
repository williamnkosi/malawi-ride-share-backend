package ServerUtils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("testing-the-new-key")

func GenerateToken(id string, phoneNumber string, firstName string, lastName string) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":              id,
		"phoneNumber":     phoneNumber,
		"firstName":       firstName,
		"lastName":        lastName,
		"exp": expirationTime,
	})
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}