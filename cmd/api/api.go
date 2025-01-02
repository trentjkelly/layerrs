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
)

type appConfig struct {
	addr string
}

type application struct {
	config 			appConfig
	trackController	*controller.TrackController
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	
	// Using middleware
	r.Use(middleware.RequestID)
  	r.Use(middleware.RealIP)
  	r.Use(middleware.Logger)
  	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 60))

	// Root route -- everything goes underneath /api
	r.Route("/api", func(r chi.Router) {
		// /api/track/
		r.Route("/track", func(r chi.Router) {
			r.Post("/", app.trackController.TrackHandlerPost)

			// r.Route("/{id}", func(r chi.Router) {
			// 	r.Get("/audio", app.trackController)
			// 	r.Get("/cover", app.trackController)
			// 	r.Get("/data", app.trackController)
			// 	r.Put("/", app.trackController)
			// 	r.Delete("/", app.trackController)
			// })
		})
	})

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