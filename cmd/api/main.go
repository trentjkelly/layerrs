package main

import (
	// Base imports
	"log"
	"fmt"
	// Loads environment variables from local .env file
	"github.com/joho/godotenv"
	// Local imports
	"github.com/trentjkelly/layerr/internals/controller"
	"github.com/trentjkelly/layerr/internals/service"
	"github.com/trentjkelly/layerr/internals/repository"
	
	"github.com/trentjkelly/layerr/internals/entities"
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

	// artistDatabaseRepository := repository.NewArtistDatabaseRepository()
	
	// artist := new(entities.Artist)
	// artist.Id = 5
	// artist.Name = "Trent Trent Trent"
	// artist.Username = "Test_3"
	// artist.Email = "test4@gmail.com"

	// err = artistDatabaseRepository.DeleteArtist(artist)
	// fmt.Println(artist)
	
	// artistDatabaseRepository.CreateArtist(artist)
	// fmt.Print(artist)

	if err != nil {
		log.Fatal(err)
	}

	// trackDatabaseRepository := repository.NewTrackDatabaseRepository()
	// trackDatabaseRepository.CreateTrack("Type Shit", "")

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
