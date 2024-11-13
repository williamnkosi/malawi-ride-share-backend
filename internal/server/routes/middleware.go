package Server

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("testing-the-new-key")

func AuthMiddleware(next http.Handler) http.Handler {
	type Claims struct {
		PhoneNumber string `json:"phoneNumber"`
		jwt.RegisteredClaims
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token is required", http.StatusBadRequest)
			return
		}

			// Assuming the token comes in the format "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization token format must be 'Bearer <token>'", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

	
		token , err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
			return jwtKey, nil
		})
	

			// Handle token parsing errors
			if err != nil {
				if errors.Is(err, jwt.ErrSignatureInvalid) {
					http.Error(w, "Invalid token signature", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			
			return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
	
			if errors.Is(err, jwt.ErrTokenExpired){
				http.Error(w, "Token has expired", http.StatusBadRequest)
			}
			next.ServeHTTP(w, r)


	})
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
