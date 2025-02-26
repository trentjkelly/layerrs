package controller

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/service"
	"fmt"
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
		fmt.Println(err)
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
	tokenString, refreshString, err := c.authService.LoginArtist(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not log in the artist", http.StatusInternalServerError)
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