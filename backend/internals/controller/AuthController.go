package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	authController := new(AuthController)
	authController.authService = authService
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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	// Validate email input
	if loginRequest.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Send login email to requested email
	c.authService.LoginArtist(r.Context(), loginRequest.Email)
	if err != nil {
		http.Error(w, "Could not send login email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Get inputs from formdata
	loginRequest := new(entities.LoginRequest)

	err := json.NewDecoder(r.Body).Decode(loginRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	// Check credentials
	tokenString, refreshString, err := c.authService.VerifyArtist(r.Context(), loginRequest.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not log in the artist", http.StatusUnauthorized)
	}

	// Send back the token string
	res := entities.LoginResponse{
		Token: tokenString,
		Refresh: refreshString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *AuthController) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	// Get refresh token
	request := new(entities.RefreshRequest)
	json.NewDecoder(r.Body).Decode(&request.RefreshToken)
	if (request.RefreshToken == "") {
		http.Error(w, "Failed to get token", http.StatusBadRequest)
	}
	log.Println(request.RefreshToken)

	// Generate new JWT
	tokenString, err := c.authService.RefreshJWT(r.Context(), request.RefreshToken)
	if err != nil {
		if err == entities.ErrInvalidToken {
			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Could not refresh jwt", http.StatusInternalServerError)
		return
	}

	// Send back refreshed jwt
	res := entities.JWTResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}