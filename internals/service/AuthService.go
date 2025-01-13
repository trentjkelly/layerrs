package service

import (
	"context"
	"fmt"
	// "log"

	"github.com/trentjkelly/layerr/internals/entities"
	"github.com/trentjkelly/layerr/internals/repository"
)

type AuthService struct {
	passwordRepository	 *repository.PasswordRepository
	artistDbRepository 	*repository.ArtistDatabaseRepository
	authRepository		*repository.AuthRepository
}

func NewAuthService(passwordRepository *repository.PasswordRepository, artistDbRepository *repository.ArtistDatabaseRepository, authRepository *repository.AuthRepository) *AuthService {
	authService := new(AuthService)
	authService.passwordRepository = passwordRepository
	authService.artistDbRepository = artistDbRepository
	authService.authRepository = authRepository
	return authService
}

// Creates a new artist with the non-optional information given
func (s *AuthService) CreateArtist(ctx context.Context, password string, username string, name string, email string) error {
	// Hash the password
	hash, err := s.passwordRepository.HashPassword(ctx, password)
	if err != nil {
		return err
	}

	// Store a new Artist using username, name, email, and hashed password
	_, err = s.artistDbRepository.CreateArtist(ctx, username, name, email, hash)
	if err != nil {
		return err
	}

	return nil
}

// Logs in an artist based on email and password
func (s *AuthService) LoginArtist(ctx context.Context, email string, password string) (string, error) {

	// Get username & password from artist
	artist := new(entities.Artist)
	err := s.artistDbRepository.GetArtistIdUsernamePassword(ctx, artist, email)
	if err != nil {
		return "", err
	}

	// Check if password is correct one
	isPassword := s.passwordRepository.CheckPassword(ctx, password, artist.Password)
	if !isPassword {
		return "", fmt.Errorf("password did not match")
	}

	// Send back a new JWT
	tokenString, err := s.authRepository.CreateJWT(artist.Id)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}