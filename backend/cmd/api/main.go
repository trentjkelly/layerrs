package main

import (
	"log"

	"github.com/trentjkelly/layerrs/internals/controller"
	"github.com/trentjkelly/layerrs/internals/service"
	"github.com/trentjkelly/layerrs/internals/repository/auth"
	"github.com/trentjkelly/layerrs/internals/repository/computing"
	"github.com/trentjkelly/layerrs/internals/repository/database"
	"github.com/trentjkelly/layerrs/internals/repository/storage"
	"github.com/trentjkelly/layerrs/internals/config"
	"github.com/joho/godotenv"

)

func main() {
	// TODO: Replace with a conditional check for the environment, when moving to docker compose
	// -- DEV ONLY --
	err := godotenv.Load(".env.backend.dev")
	if err != nil {
		log.Fatal(err)
	}
	// -- END DEV ONLY --

	// Database Connection
	pool, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// -- REPOSITORIES --
	// Computing Repositories
	trackConversionRepo := computingRepository.NewTrackConversionRepository()
	waveformRepo := computingRepository.NewWaveformHeightsRepository()

	// Auth Repositories
	passwordRepo := authRepository.NewPasswordRepository()
	authRepo := authRepository.NewAuthRepository()

	// Database Repositories
	artistDatabaseRepo := databaseRepository.NewArtistDatabaseRepository(pool)
	likesDatabaseRepo := databaseRepository.NewLikesDatabaseRepository(pool)
	trackDatabaseRepo := databaseRepository.NewTrackDatabaseRepository(pool)
	trackTreeDatabaseRepo := databaseRepository.NewTrackTreeDatabaseRepository(pool)

	// Storage Repositories
	coverStorageRepo := storageRepository.NewCoverStorageRepository()
	// portraitStorageRepo := repository.NewPortraitStorageRepository()
	trackStorageRepo := storageRepository.NewTrackStorageRepository()

	// -- SERVICES --
	authService := service.NewAuthService(passwordRepo, artistDatabaseRepo, authRepo)
	trackService := service.NewTrackService(trackStorageRepo, coverStorageRepo, trackDatabaseRepo, trackTreeDatabaseRepo, trackConversionRepo, waveformRepo)
	recService := service.NewRecommendationsService(trackDatabaseRepo, likesDatabaseRepo)
	artistService := service.NewArtistService(artistDatabaseRepo)
	likesService := service.NewLikesService(likesDatabaseRepo, trackDatabaseRepo)

	// -- CONTROLLERS --
	authController := controller.NewAuthController(authService)
	trackController := controller.NewTrackController(trackService)
	recController := controller.NewRecommendationsController(recService)
	likesController := controller.NewLikesController(likesService)
	artistController := controller.NewArtistController(artistService)

	// -- CONFIGURATION --
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
