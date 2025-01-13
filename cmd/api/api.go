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
	authController 				*controller.AuthController
	trackController				*controller.TrackController
	recommendationsController 	*controller.RecommendationsController
	likesController				*controller.LikesController
	artistController 			*controller.ArtistController
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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Root route -- everything goes underneath /api
	r.Route("/api", func(r chi.Router) {

		r.Route("/authentication", func(r chi.Router) {
			// r.Options("/login", app.trackController.AuthHandlerOptions)
			r.Post("/signup", app.authController.RegisterArtistHandler)
			r.Post("/login", app.authController.LogInArtistHandler)
			r.Post("/refresh", app.authController.RefreshHandler)
		})

		r.Route("/artist", func(r chi.Router) {
			r.Get("/{artistId}", app.artistController.ArtistHandlerGet)
		})

		r.Route("/track", func(r chi.Router) {
			r.Options("/", app.trackController.TrackHandlerOptions)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/audio", app.trackController.TrackAudioHandlerGet)
				r.Get("/cover", app.trackController.TrackCoverHandlerGet)
				r.Get("/data", app.trackController.TrackerDataHandlerGet)

				// r.Use(AuthJWTMiddleware)
				// r.Put("/", app.trackController.)
				// r.Delete("/", app.trackController)
			})
			r.Group(func (r chi.Router) {
				r.Use(AuthJWTMiddleware)
				r.Post("/", app.trackController.TrackHandlerPost)
			})
		})

		// Different algorithms for showing pages of songs
		r.Route("/recommendations", func(r chi.Router) {
			r.Route("/home", func (r chi.Router) {
				r.Get("/", app.recommendationsController.RecommendationsHandlerHomeGet) // Base home page algorithm

				r.Group(func (r chi.Router) {
					r.Use(AuthJWTMiddleware)
					r.Get("/{artistId}", app.recommendationsController.RecommendationsHandlerHomeGet) // Personalized home page algorithm
				})
			})
			r.Route("/library", func (r chi.Router) {
				r.Use(AuthJWTMiddleware)
				r.Get("/{artistId}", app.recommendationsController.ReccomendationsHandlerLikesGet) // User's library algorithm
			})
		})

		r.Route("/likes", func(r chi.Router) {
			r.Use(AuthJWTMiddleware)
			r.Options("/", app.likesController.LikesHandlerOptions)
			r.Post("/", app.likesController.LikesHandlerPost)
			r.Get("/", app.likesController.LikesHandlerGet)
			r.Delete("/", app.likesController.LikesHandlerDelete)
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