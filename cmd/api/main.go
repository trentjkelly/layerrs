package main

import (
	"log"
)

func main() {
	// Setup configuration and injected dependencies
	cfg := config{
		addr : ":8080",
	}

	app := &application{
		config: cfg,
	}

	// Mount and run the application
	mux := app.mount()
	log.Fatal(app.run(mux))
}