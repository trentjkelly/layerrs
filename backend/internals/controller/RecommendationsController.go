package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/service"
)

type RecommendationsController struct {
	recService *service.RecommendationsService
}

func NewRecommendationsController(recService *service.RecommendationsService) *RecommendationsController {
	recController := new(RecommendationsController)
	recController.recService = recService
	return recController
}

// Sends a user what tracks to show on their homepage
func (c *RecommendationsController) RecommendationsHandlerHomeGet(w http.ResponseWriter, r *http.Request) {
	artistId := 0
	headerString := r.Header.Get("Authorization")
	if headerString != "" {
		parts := strings.Split(headerString, " ")
		if len(parts) == 2 {
			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("AUTH_SECRET_KEY")), nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if sub, ok := claims["sub"].(float64); ok {
						artistId = int(sub)
					}
				}
			}
		}
	}

	rec, err := c.recService.MostRecentAlgorithm(r.Context(), artistId)
	if err != nil {
		log.Printf("[ERROR] RecommendationsHandlerHomeGet: %s", err)
		http.Error(w, "Unable to get reccomendations", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(rec)
	if err != nil {
		log.Printf("[ERROR] RecommendationsHandlerHomeGet: %s", err)
		http.Error(w, "Unable to encode recommendations to json", http.StatusInternalServerError)
		return
	}
}

// Sends a user what tracks to show on their likes page
func (c *RecommendationsController) RecommendationsHandlerLibraryGet(w http.ResponseWriter, r *http.Request) {
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistId := int(artistIdFloat)

	likesArr, err := c.recService.ArtistLikesAlgorithm(r.Context(), artistId, 0)
	if err != nil {
		log.Printf("[ERROR] RecommendationsHandlerLibraryGet: %s", err)
		http.Error(w, "Could not retrieve liked tracks", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(likesArr)
	if err != nil {
		log.Printf("[ERROR] RecommendationsHandlerLibraryGet: %s", err)
		http.Error(w, "Unable to encode recommendations to json", http.StatusInternalServerError)
		return
	}
}
