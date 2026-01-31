package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"

	"log"
	"fmt"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/repository/auth"
	"github.com/trentjkelly/layerrs/internals/repository/database"
)

type AuthService struct {
	passwordRepository	 *authRepository.PasswordRepository
	artistDbRepository 	*databaseRepository.ArtistDatabaseRepository
	authRepository		*authRepository.AuthRepository
	verificationEmailRepository *authRepository.VerificationEmailRepository
}

func NewAuthService(passwordRepository *authRepository.PasswordRepository, artistDbRepository *databaseRepository.ArtistDatabaseRepository, authRepository *authRepository.AuthRepository, verificationEmailRepository *authRepository.VerificationEmailRepository) *AuthService {
	authService := new(AuthService)
	authService.passwordRepository = passwordRepository
	authService.artistDbRepository = artistDbRepository
	authService.authRepository = authRepository
	authService.verificationEmailRepository = verificationEmailRepository
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

	// Send a verification email to the user
	// err = s.verificationEmailRepository.SendVerificationEmail(email)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Logs in an artist through magic links
func (s *AuthService) LoginArtist(ctx context.Context, email string) error {

	// TODO: Generate a new magic link token


	// Create a new Verification Email
	verificationEmail := entities.LoginEmailInfo{
		EmailSender: "Layerrs <team@login.layerrs.com>",
		EmailRecipients: []string{email},
		EmailSubject: "Login to Layerrs",
		EmailBodyHTML: "<p>Please click the link below to login to your Layerrs account.</p>",
	}

	// Send a login email to the user
	err := s.verificationEmailRepository.SendEmail(verificationEmail)
	if err != nil {
		return fmt.Errorf("could not send login email: %w", err)
	}

	return nil
}

// Logs in an artist based on email magic link verification
func (s *AuthService) VerifyArtist(ctx context.Context, email string) (string, string, error) {
	// // Get a new JWT
	// tokenString, err := s.authRepository.CreateJWT(artist.Id)
	// if err != nil {
	// 	return "", "", err
	// }

	// // Get a new refresh token
	// refreshString, err := s.authRepository.CreateRefreshToken(artist.Id)
	// if err != nil {
	// 	return "", "", err
	// }

	// return tokenString, refreshString, nil
	return "", "", nil
}

func (s *AuthService) RefreshJWT(ctx context.Context, refreshToken string) (string, error) {
	// Check if the refresh token is valid still
	log.Println(refreshToken)
	token, err := s.authRepository.ValidateJWT(ctx, refreshToken)
	if err != nil {
		return "", entities.ErrInvalidToken
	}

	// Get claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("here2")
		return "", err
	}
	// Get artistId from claims
	artistInterface, ok := claims["sub"]
	if !ok {
		log.Println("here3")
		return "", err
	}
	artistId := artistInterface.(int)

	// Create new JWT
	jwt, err := s.authRepository.CreateJWT(artistId)
	if err != nil {
		log.Println("here4")
		return "", err
	}

	return jwt, nil
}