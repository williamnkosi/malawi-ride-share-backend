package Server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

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


type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Recovery middleware to handle panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)

				// Respond with a 500 Internal Server Error
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := Response{
					Message: "Internal Server Error",
					Status:  http.StatusInternalServerError,
				}
				json.NewEncoder(w).Encode(response)
			}
		}()

		next.ServeHTTP(w, r)
	})
}