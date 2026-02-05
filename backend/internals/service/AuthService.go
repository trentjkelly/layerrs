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

const MAX_RETRIES = 10

type AuthService struct {
	passwordRepository	 *authRepository.PasswordRepository
	artistDbRepository 	*databaseRepository.ArtistDatabaseRepository
	authRepository		*authRepository.AuthRepository
	verificationEmailRepository *authRepository.VerificationEmailRepository
	authDatabaseRepository *databaseRepository.AuthDatabaseRepository
	magicLink string
}

func NewAuthService(passwordRepository *authRepository.PasswordRepository, artistDbRepository *databaseRepository.ArtistDatabaseRepository, authRepository *authRepository.AuthRepository, verificationEmailRepository *authRepository.VerificationEmailRepository, authDatabaseRepository *databaseRepository.AuthDatabaseRepository, magicLink string) *AuthService {
	authService := new(AuthService)
	authService.passwordRepository = passwordRepository
	authService.artistDbRepository = artistDbRepository
	authService.authRepository = authRepository
	authService.verificationEmailRepository = verificationEmailRepository
	authService.authDatabaseRepository = authDatabaseRepository
	authService.magicLink = magicLink
	return authService
}

// Logs in an artist through magic links
func (s *AuthService) LoginArtist(ctx context.Context, email string) error {
	var artistId int

	// Get the artist from the database / see if they exist
	artist, err := s.artistDbRepository.GetArtistByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("could not get artist from database: %w", err)
	}

	if artist == nil {
		artistId, err = s.createArtistLoop(ctx, email)
		if err != nil {
			return fmt.Errorf("could not create artist in database: %w", err)
		}
	} else {
		artistId = artist.Id
	}

	// Generate a new magic link token
	magicLinkToken, err := s.authRepository.CreateMagicLinkToken()
	if err != nil {
		return fmt.Errorf("could not create magic link token: %w", err)
	}

	// Store the magic link token in the database
	err = s.authDatabaseRepository.CreateMagicLinkToken(ctx, magicLinkToken, artistId)
	if err != nil {
		return fmt.Errorf("could not store magic link token in database: %w", err)
	}

	verifyPath := fmt.Sprintf("%s%s", s.magicLink, magicLinkToken)

	verificationEmail := entities.LoginEmailInfo{
		EmailSender: "Layerrs <team@login.layerrs.com>",
		EmailRecipients: []string{email},
		EmailSubject: "Login to Layerrs",
		EmailBodyHTML: fmt.Sprintf("<p>Please click the link below to login to your Layerrs account.</p><a href=\"%s\">Login to Layerrs</a>", verifyPath),
	}

	// Send a login email to the user
	err = s.verificationEmailRepository.SendEmail(verificationEmail)
	if err != nil {
		return fmt.Errorf("could not send login email: %w", err)
	}

	return nil
}

func (s *AuthService) createArtistLoop(ctx context.Context, email string) (int, error) {
	var artistId int
	var err error

	for i := 0; i < MAX_RETRIES; i++ {
		username, usernameErr := s.authRepository.CreateRandomUsername()
		if usernameErr != nil {
			return 0, fmt.Errorf("could not create random username: %w", usernameErr)
		}
		
		artistId, err = s.artistDbRepository.CreateArtist(ctx, username, email)
		if err == nil {
			break
		}
	}

	if err != nil {
		return 0, fmt.Errorf("could not create artist in database: %w", err)
	}

	return artistId, nil
}

// TODO: Verify should check, and return a redirect link to the frontend
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