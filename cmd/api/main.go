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
	// artistDatabaseRepo := repository.NewArtistDatabaseRepository()
	// likesDatabaseRepo := repository.NewLikesDatabaseRepository()
	trackDatabaseRepo := repository.NewTrackDatabaseRepository()
	trackTreeDatabaseRepo := repository.NewTrackTreeDatabaseRepository()
	coverStorageRepo := repository.NewCoverStorageRepository()
	// portraitStorageRepo := repository.NewPortraitStorageRepository()
	trackStorageRepo := repository.NewTrackStorageRepository()

	// Services
	trackService := service.NewTrackService(trackStorageRepo, coverStorageRepo, trackDatabaseRepo, trackTreeDatabaseRepo)
	// artistService := service.NewArtistService()
	// likesService := service.NewLikesService()

	// Controllers
	trackController := controller.NewTrackController(trackService)
	// artistController := controller.NewArtistController(artistService)
	// likesController := controller.NewLikesController(likesService)

	
	// Setup configuration and injected dependencies
	cfg := appConfig{
		addr : ":8080",
	}

	app := &application{
		config		: cfg,
		trackController: trackController,
	}

	// Mount and run the application
	mux := app.mount()
	log.Fatal(app.run(mux))
}
