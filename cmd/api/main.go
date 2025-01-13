package main

import (
	// Base imports
	"log"
	// Loads environment variables from local .env file
	"github.com/joho/godotenv"
	// Local imports
	"github.com/trentjkelly/layerr/internals/controller"
	"github.com/trentjkelly/layerr/internals/service"
	"github.com/trentjkelly/layerr/internals/repository"
)

func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// 
	// Initialize Layered Architecture
	// 

	// Repositories
	passwordRepo := repository.NewPasswordRepository()
	authRepo := repository.NewAuthRepository()
	artistDatabaseRepo := repository.NewArtistDatabaseRepository()
	likesDatabaseRepo := repository.NewLikesDatabaseRepository()
	trackDatabaseRepo := repository.NewTrackDatabaseRepository()
	trackTreeDatabaseRepo := repository.NewTrackTreeDatabaseRepository()
	coverStorageRepo := repository.NewCoverStorageRepository()
	// portraitStorageRepo := repository.NewPortraitStorageRepository()
	trackStorageRepo := repository.NewTrackStorageRepository()

	// Services
	authService := service.NewAuthService(passwordRepo, artistDatabaseRepo, authRepo)
	trackService := service.NewTrackService(trackStorageRepo, coverStorageRepo, trackDatabaseRepo, trackTreeDatabaseRepo)
	recService := service.NewRecommendationsService(trackDatabaseRepo)
	artistService := service.NewArtistService(artistDatabaseRepo)
	likesService := service.NewLikesService(likesDatabaseRepo, trackDatabaseRepo)

	// Controllers
	authController := controller.NewAuthController(authService)
	trackController := controller.NewTrackController(trackService)
	recController := controller.NewRecommendationsController(recService)
	likesController := controller.NewLikesController(likesService)
	artistController := controller.NewArtistController(artistService)

	
	// Setup configuration and injected dependencies
	cfg := appConfig{
		addr : ":8080",
	}

	app := &application{
		config		: cfg,
		trackController: trackController,
		recommendationsController: recController,
		authController: authController,
		likesController: likesController,
		artistController: artistController,
	}

	// Mount and run the application
	mux := app.mount()
	log.Fatal(app.run(mux))
}
