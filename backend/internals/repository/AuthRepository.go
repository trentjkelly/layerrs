package repository

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
	"fmt"
	"context"
)

type AuthRepository struct {
	secretKey	string
}

func NewAuthRepository() *AuthRepository {
	authRepo := new(AuthRepository)
	authRepo.secretKey = os.Getenv("AUTH_SECRET_KEY")
	return authRepo
}

// Creates a new JWT for a logged in user
func (r *AuthRepository) CreateJWT(artistId int) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": artistId, // subject
		"iss": "layerr", // issuer
		"aud": "artist", // audience (role)
		"exp": time.Now().Add(time.Hour).Unix(), //expiration
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString([]byte(r.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Create a new refresh token
func (r *AuthRepository) CreateRefreshToken(artistId int) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": artistId, // subject
		"iss": "layerr", // issuer
		"aud": "artist", // audience (role)
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), //expiration
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString([]byte(r.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Check if a given JWT is valid
func (r *AuthRepository) ValidateJWT(ctx context.Context, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	return token, nil
}
