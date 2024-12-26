package main

import (
	// Native packages
	"log"
	"net/http"
	"time"
	// Chi router
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// Local packages
	"github.com/trentjkelly/layerr/internals/controller"
	// "github.com/trentjkelly/layerr/internals/service"
	"github.com/trentjkelly/layerr/internals/repository"

)

type appConfig struct {
	addr string
}

type application struct {
	config 		appConfig
	trackStorageRepository	*repository.TrackStorageRepository
	trackController *controller.TrackController
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	
	// Using middleware
	r.Use(middleware.RequestID)
  	r.Use(middleware.RealIP)
  	r.Use(middleware.Logger)
  	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 60))

	r.Get("/track",  app.trackController.TrackHandlerGet)
	r.Put("/track", app.trackController.TrackHandlerPut)
	r.Options("/track", app.trackController.TrackHandlerOptions)
	return r
}

func (app *application) run(mux http.Handler) error {	
	server := &http.Server{
		Addr:			app.config.addr,
		Handler:		mux,
		WriteTimeout:	time.Second * 30,
		ReadTimeout: 	time.Second * 10,
		IdleTimeout: 	time.Minute,
	}

	log.Printf("Server has started at %s", server.Addr)

	return server.ListenAndServe()
}