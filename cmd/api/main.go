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

	// Initialize Layered Architecture
	trackStorageRepository := repository.NewTrackStorageRepository()
	trackService := service.NewTrackService(trackStorageRepository)
	trackController := controller.NewTrackController(trackService)

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
