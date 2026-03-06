package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/service"
)

type AuthController struct {
	authService *service.AuthService
	frontendUrl string
}

func NewAuthController(authService *service.AuthService, frontendUrl string) *AuthController {
	authController := new(AuthController)
	authController.authService = authService
	authController.frontendUrl = frontendUrl
	return authController
}

// OPTIONS request for browsers when they test for CORS before PUT request
func (c *TrackController) AuthHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusNoContent)
}

func (c *AuthController) LoginArtistHandler(w http.ResponseWriter, r *http.Request) {
	// Get inputs from the formdata
	loginRequest := new(entities.LoginRequest)
	err := json.NewDecoder(r.Body).Decode(loginRequest)
	if err != nil {
		log.Println("Invalid JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate email input
	if loginRequest.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Send login email to requested email
	err = c.authService.LoginArtist(r.Context(), loginRequest.Email)
	if err != nil {
		log.Println("Could not send login email:", err)
		http.Error(w, "Could not send login email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Get inputs from formdata
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Check credentials
	tokenString, refreshString, isFirstLogin, err := c.authService.VerifyArtist(r.Context(), token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not log in the artist", http.StatusUnauthorized)
		return
	}

	redirectURL := fmt.Sprintf("%s/login/callback#jwt=%s&refresh=%s", c.frontendUrl, url.QueryEscape(tokenString), url.QueryEscape(refreshString))
	if isFirstLogin {
		redirectURL += "&firstLogin=true"
	}
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (c *AuthController) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	request := new(entities.RefreshRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		log.Println("Invalid JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if request.RefreshToken == "" {
		log.Println("Failed to get token")
		http.Error(w, "Failed to get token", http.StatusBadRequest)
		return
	}

	// Generate new JWT
	tokenString, err := c.authService.RefreshJWT(r.Context(), request.RefreshToken)
	if err != nil {
		if err == entities.ErrInvalidToken {
			log.Println("Token is invalid")
			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}
		log.Println("refresh error", err)
		http.Error(w, "Could not refresh jwt", http.StatusInternalServerError)
		return
	}

	// Send back refreshed jwt
	res := entities.JWTResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}