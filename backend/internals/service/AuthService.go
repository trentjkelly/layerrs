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
		artistId, err = s.artistDbRepository.CreateArtist(ctx, email)
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

// Logs in an artist based on email magic link verification
func (s *AuthService) VerifyArtist(ctx context.Context, token string) (string, string, bool, error) {
	magicLinkToken, err := s.authDatabaseRepository.GetMagicLinkToken(ctx, token)
	if err != nil {
		return "", "", false, fmt.Errorf("could not get magic link token from database: %w", err)
	}

	isValid, err := s.authRepository.VerifyMagicLinkToken(ctx, magicLinkToken)
	if err != nil || !isValid {
		return "", "", false, fmt.Errorf("could not verify magic link token: %w", err)
	}

	jwt, err := s.authRepository.CreateJWT(magicLinkToken.ArtistId)
	if err != nil {
		return "", "", false, fmt.Errorf("could not create JWT: %w", err)
	}

	refreshToken, err := s.authRepository.CreateRefreshToken(magicLinkToken.ArtistId)
	if err != nil {
		return "", "", false, fmt.Errorf("could not create refresh token: %w", err)
	}

	count, err := s.authDatabaseRepository.CountMagicLinkTokensByArtistId(ctx, magicLinkToken.ArtistId)
	if err != nil {
		log.Printf("[WARN] VerifyArtist: could not count magic link tokens: %s", err)
		return jwt, refreshToken, false, nil
	}

	return jwt, refreshToken, count == 1, nil
}

// Refresh a JWT token using a refresh token
func (s *AuthService) RefreshJWT(ctx context.Context, refreshToken string) (string, error) {
	log.Println(refreshToken)
	token, err := s.authRepository.ValidateJWT(ctx, refreshToken)
	if err != nil {
		return "", entities.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	sub, ok := claims["sub"]
	if !ok {
		return "", err
	}

	artistIdFloat, ok := sub.(float64)
	if !ok {
		return "", entities.ErrInvalidToken
	}
	artistId := int(artistIdFloat)

	jwt, err := s.authRepository.CreateJWT(artistId)
	if err != nil {
		log.Println("here4")
		return "", err
	}

	return jwt, nil
}