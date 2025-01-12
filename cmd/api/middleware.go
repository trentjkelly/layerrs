package main

import (
	"net/http"
	"os"
	"strings"
)

func AuthJWTMiddleware(next http.Handler) http.Handler {

	// Load the secret key once
	secretKey := os.Getenv("AUTH_SECRET_KEY")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// Extract string from auth header
		headerString := r.Header.Get("Authorization")
		if headerString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Format should be Bearer aks7f9shfhsd...
		authString := strings.Split(headerString, " ")
		if len(authString) != 2 {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Get the actual token string & secret key
		tokenString := authString[1]
		_, err := ValidateJWT(tokenString, secretKey)
		if err != nil {
			http.Error(w, "Invalid JWT", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}