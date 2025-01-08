package main

import (
	// Native packages
	"log"
	"net/http"
	"time"
	// Chi router
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi/v5/middleware"
	// Local packages
	"github.com/trentjkelly/layerr/internals/controller"
)

type appConfig struct {
	addr string
}

type application struct {
	config 						appConfig
	trackController				*controller.TrackController
	recommendationsController 	*controller.RecommendationsController
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	
	// Using middleware
	r.Use(middleware.RequestID)
  	r.Use(middleware.RealIP)
  	r.Use(middleware.Logger)
  	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 60))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Root route -- everything goes underneath /api
	r.Route("/api", func(r chi.Router) {

		r.Route("/track", func(r chi.Router) {
			r.Options("/", app.trackController.TrackHandlerOptions)
			r.Post("/", app.trackController.TrackHandlerPost)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/audio", app.trackController.TrackAudioHandlerGet)
				r.Get("/cover", app.trackController.TrackCoverHandlerGet)
				r.Get("/data", app.trackController.TrackerDataHandlerGet)
				// r.Put("/", app.trackController)
				// r.Delete("/", app.trackController)
			})
		})

		r.Route("/recommendations", func(r chi.Router) {
			r.Route("/home", func (r chi.Router){
				r.Get("/{artistId}", app.recommendationsController.RecommendationsHandlerHomeGet)
			})
		})

		// r.Route("/likes", func(r chi.Router) {
		// 	r.Post("/", app.LikesController.LikesHandlerPost)
		// 	r.Get("/", app.LikesController.LikesHandlerGet)
		// 	r.Delete("/{id}", app.LikesController.LikesHandlerDelete)
		// })

		// r.Route("/artist", func(r chi.Router) {})
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