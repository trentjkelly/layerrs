package authRepository

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
	"fmt"
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/trentjkelly/layerrs/internals/entities"
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
		"exp": time.Now().Add(time.Minute * 15).Unix(), //expiration
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
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // expiration
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

func (r *AuthRepository) CreateMagicLinkToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("could not read random bytes for magic link token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (r *AuthRepository) CreateRandomUsername() (string, error) {
	usernameLength := 64
	prefix := "user"

	prefixLength := len(prefix)

	
	b := make([]byte, usernameLength - prefixLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("could not read random bytes for username: %w", err)
	}

	username := prefix + string(b)
	return username, nil
}

// Check if the magic link token is valid (not expired)
func (r *AuthRepository) VerifyMagicLinkToken(ctx context.Context, token entities.MagicLinkToken) (bool, error) {
	if token.CreatedAt.Before(time.Now().Add(-time.Minute * 15)) {
		return false, fmt.Errorf("magic link token has expired")
	}

	return true, nil
}