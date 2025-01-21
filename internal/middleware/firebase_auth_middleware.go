package Middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseAuthMiddleware verifies Firebase ID tokens.
func FirebaseAuthMiddleware(next http.Handler) http.Handler {

	// Initialize Firebase App
	opt := option.WithCredentialsFile("/Users/williamnkosi/repo/malawi-ride-share-backend/cmd/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Firebase App: %v", err))
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Failed to create Firebase Auth client: %v", err))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Check if the header is in the format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		idToken := parts[1]
		// Verify the ID token
		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to verify ID token: %v", err), http.StatusUnauthorized)
			return
		}
		print("verified")

		// Add token claims to the request context
		ctx := context.WithValue(r.Context(), "firebaseUser", token)
		print("firebase----Done")
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
