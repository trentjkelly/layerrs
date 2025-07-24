package main

import (
	"log"
	"os"

	"github.com/trentjkelly/layerrs/internals/config"
	"github.com/trentjkelly/layerrs/internals/controller"
	"github.com/trentjkelly/layerrs/internals/repository/auth"
	"github.com/trentjkelly/layerrs/internals/repository/computing"
	"github.com/trentjkelly/layerrs/internals/repository/database"
	"github.com/trentjkelly/layerrs/internals/repository/storage"
	"github.com/trentjkelly/layerrs/internals/service"
)

const (
	ENVIRONMENT = "ENV"
	DEVELOPMENT = "DEVELOPMENT"
	PRODUCTION = "PRODUCTION"
)

func main() {
	// Get the environment
	env := os.Getenv(ENVIRONMENT)
	if env == "" {
		log.Fatal("Could not find the environment variable ENV")
	}

	// Database Connection
	pool, err := config.InitDB(env)
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
	coverStorageRepo := storageRepository.NewCoverStorageRepository(env)
	// portraitStorageRepo := storageRepository.NewPortraitStorageRepository(env)
	trackStorageRepo := storageRepository.NewTrackStorageRepository(env)

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
