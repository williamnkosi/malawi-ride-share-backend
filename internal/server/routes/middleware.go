package Server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("testing-the-new-key")


func AuthMiddleware() gin.HandlerFunc {
	type Claims struct {
		PhoneNumber string `json:"phoneNumber"`
		jwt.RegisteredClaims
	}

	return func(c *gin.Context) {
		print("entered")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// Assuming the token comes in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token format must be 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Handle token parsing errors
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Ensure the token isn't expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		// Set the username in the context for the next handlers
		fmt.Println("======]")
		print(claims.PhoneNumber)

		c.Set("phoneNumber", claims.PhoneNumber)
		c.Next()

	}
}

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error (you could also send it to an error tracking service)
				log.Printf("Panic recovered: %s", err)

				// Return a 500 status with a custom message
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error. Please try again later.",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
