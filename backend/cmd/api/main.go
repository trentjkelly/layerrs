package main

import (
	// Base imports
	"log"
	// Loads environment variables from local .env file
	"github.com/trentjkelly/layerrs/internals/controller"
	"github.com/trentjkelly/layerrs/internals/service"
	"github.com/trentjkelly/layerrs/internals/repository/auth"
	"github.com/trentjkelly/layerrs/internals/repository/computing"
	"github.com/trentjkelly/layerrs/internals/repository/database"
	"github.com/trentjkelly/layerrs/internals/repository/storage"

	// -- DEV ONLY --
	// "github.com/joho/godotenv"
	// -- END DEV ONLY --
)

func main() {

	// -- DEV ONLY --
	// err := godotenv.Load(".env.backend.dev")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// -- END DEV ONLY --
	
	// 
	// Initialize Layered Architecture
	// 

	// Repositories
	trackConversionRepo := computingRepository.NewTrackConversionRepository()
	waveformRepo := computingRepository.NewWaveformHeightsRepository()
	passwordRepo := authRepository.NewPasswordRepository()
	authRepo := authRepository.NewAuthRepository()
	artistDatabaseRepo := databaseRepository.NewArtistDatabaseRepository()
	likesDatabaseRepo := databaseRepository.NewLikesDatabaseRepository()
	trackDatabaseRepo := databaseRepository.NewTrackDatabaseRepository()
	trackTreeDatabaseRepo := databaseRepository.NewTrackTreeDatabaseRepository()
	coverStorageRepo := storageRepository.NewCoverStorageRepository()
	// portraitStorageRepo := repository.NewPortraitStorageRepository()
	trackStorageRepo := storageRepository.NewTrackStorageRepository()

	// Services
	authService := service.NewAuthService(passwordRepo, artistDatabaseRepo, authRepo)
	trackService := service.NewTrackService(trackStorageRepo, coverStorageRepo, trackDatabaseRepo, trackTreeDatabaseRepo, trackConversionRepo, waveformRepo)
	recService := service.NewRecommendationsService(trackDatabaseRepo, likesDatabaseRepo)
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
