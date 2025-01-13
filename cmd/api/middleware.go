package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/trentjkelly/layerr/internals/entities"
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
		token, err := ValidateJWT(tokenString, secretKey)
		if err != nil {
			http.Error(w, "Invalid JWT", http.StatusUnauthorized)
			return
		}

		// Get the subject id from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Could not get the claims from JWT", http.StatusUnauthorized)
			return
		}

		// Get the sub from the claims
		artistId, ok := claims["sub"]
		if !ok {
			http.Error(w, "Invalid subject claim", http.StatusUnauthorized)
			return
		}

		// Pass artistId to next handler as context
		ctx := context.WithValue(r.Context(), entities.ArtistIdKey, artistId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
