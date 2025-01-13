package repository

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
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
