package controller

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"

	"github.com/trentjkelly/layerr/internals/entities"
	"github.com/trentjkelly/layerr/internals/service"
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

func (c *AuthController) RegisterArtistHandler(w http.ResponseWriter, r *http.Request) {
	
	// Get inputs from the formdata
	signupRequest := new(entities.SignupRequest)
	err := json.NewDecoder(r.Body).Decode(signupRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	// Create new artist
	err = c.authService.CreateArtist(r.Context(), signupRequest.Password, signupRequest.Username, signupRequest.Name, signupRequest.Email)
	if err != nil {
		http.Error(w, "Could not create new artst", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) LogInArtistHandler(w http.ResponseWriter, r *http.Request) {

	// Get inputs from formdata
	loginRequest := new(entities.LoginRequest)

	err := json.NewDecoder(r.Body).Decode(loginRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	// Check credentials
	tokenString, err := c.authService.LoginArtist(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not log in the artist", http.StatusInternalServerError)
	}

	// Send back the token string
	res := entities.LoginResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}