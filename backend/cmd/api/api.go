package main

import (
	// Native packages
	"log"
	"net/http"
	"time"
	// Chi router
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	// Local packages
	"github.com/trentjkelly/layerrs/internals/controller"
)

type appConfig struct {
	addr string
}

type application struct {
	config                    appConfig
	authController            *controller.AuthController
	trackController           *controller.TrackController
	recommendationsController *controller.RecommendationsController
	likesController           *controller.LikesController
	artistController          *controller.ArtistController
	layerrsController         *controller.LayerrsController
	sellerController          *controller.SellerController
	webhookController         *controller.WebhookController
	purchaseController        *controller.PurchaseController
	pageController            *controller.PageController
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
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Root route -- everything goes underneath /api
	r.Route("/api", func(r chi.Router) {

		r.Route("/authentication", func(r chi.Router) {
			// r.Options("/login", app.trackController.AuthHandlerOptions)
			r.Post("/login", app.authController.LoginArtistHandler)
			r.Route("/verify", func(r chi.Router) {
				r.Get("/", app.authController.VerifyEmailHandler)
			})
			r.Post("/refresh", app.authController.RefreshHandler)
		})

		r.Route("/profile", func(r chi.Router) {
			r.Use(AuthJWTMiddleware)
			r.Put("/", app.artistController.ArtistHandlerPut)
			r.Get("/", app.artistController.ArtistHandlerGet)
		})

		// r.Route("/artist", func(r chi.Router) {
		// 	r.Get("/{artistId}", app.artistController.ArtistHandlerGet)
		// })

		r.Route("/track", func(r chi.Router) {
			r.Options("/", app.trackController.TrackHandlerOptions)
			r.Post("/batch", app.trackController.TrackBatchHandlerPost)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/audio", app.trackController.TrackAudioHandlerGet)
			r.Group(func(r chi.Router) {
				r.Use(AuthJWTMiddleware)
				r.Get("/download", app.trackController.TrackDownloadHandlerGet)
				r.Put("/price", app.trackController.UpdateTrackPriceHandlerPut)
				r.Post("/play", app.trackController.TrackPlayHandlerPost)
			})
			r.Get("/data", app.trackController.TrackerDataHandlerGet)
			r.Get("/recommendation", app.trackController.TrackRecommendationHandlerGet)
			r.Get("/graph", app.trackController.TrackGraphHandlerGet)
		})
			r.Group(func(r chi.Router) {
				r.Use(AuthJWTMiddleware)
				r.Post("/", app.trackController.TrackHandlerPost)
			})
		})

		// Different algorithms for showing pages of songs
		r.Route("/recommendations", func(r chi.Router) {
			r.Route("/home", func(r chi.Router) {
				r.Get("/", app.recommendationsController.RecommendationsHandlerHomeGet) // Base home page algorithm

				r.Group(func(r chi.Router) {
					r.Use(AuthJWTMiddleware)
					r.Get("/{artistId}", app.recommendationsController.RecommendationsHandlerHomeGet) // Personalized home page algorithm
				})
			})

			r.Route("/library", func(r chi.Router) {
				r.Route("/likes", func(r chi.Router) {
					r.Use(AuthJWTMiddleware)
					r.Get("/", app.recommendationsController.RecommendationsHandlerLibraryLikesGet) // User's liked tracks
				})
			})
		})

		r.Route("/likes", func(r chi.Router) {
			r.Use(AuthJWTMiddleware)
			r.Options("/", app.likesController.LikesHandlerOptions)
			r.Post("/", app.likesController.LikesHandlerPost)
			r.Get("/", app.likesController.LikesHandlerGet)
			r.Delete("/", app.likesController.LikesHandlerDelete)
		})

		r.Route("/layerrs", func(r chi.Router) {
			r.Use(AuthJWTMiddleware)
			r.Options("/", app.layerrsController.LayerrsHandlerOptions)
			r.Get("/", app.layerrsController.LayerrsHandlerGet)
		})

		r.Route("/seller", func(r chi.Router) {
			r.Route("/onboard", func(r chi.Router) {
				r.Use(AuthJWTMiddleware)
				r.Options("/", app.sellerController.OnboardSellerHandlerOptions)
				r.Post("/", app.sellerController.OnboardSellerHandlerPost)
			})
		})

		r.Route("/purchase", func(r chi.Router) {
			r.Use(AuthJWTMiddleware)
			r.Options("/checkout", app.purchaseController.CheckoutHandlerOptions)
			r.Post("/checkout", app.purchaseController.CheckoutHandlerPost)
			r.Options("/confirm", app.purchaseController.ConfirmHandlerOptions)
			r.Post("/confirm", app.purchaseController.ConfirmHandlerPost)
			r.Options("/my-purchases", app.purchaseController.GetPurchasesHandlerOptions)
			r.Get("/my-purchases", app.purchaseController.GetPurchasesHandlerGet)
			r.Route("/{trackId}", func(r chi.Router) {
				r.Options("/download", app.purchaseController.GetDownloadURLHandlerOptions)
				r.Get("/download", app.purchaseController.GetDownloadURLHandlerGet)
			})
		})

		r.Route("/pages", func(r chi.Router) {
			r.With(AuthJWTMiddleware).Post("/", app.pageController.CreatePageHandler)
		r.Route("/{pageId}", func(r chi.Router) {
			r.With(OptionalAuthJWTMiddleware).Get("/", app.pageController.GetPageHandler)
			r.Group(func(r chi.Router) {
				r.Use(AuthJWTMiddleware)
					r.Patch("/", app.pageController.UpdatePageHandler)
					r.Delete("/", app.pageController.DeletePageHandler)
					r.Post("/follow", app.pageController.FollowPageHandler)
					r.Delete("/follow", app.pageController.UnfollowPageHandler)
					r.Route("/tracks", func(r chi.Router) {
						r.Post("/", app.pageController.AddTrackToPageHandler)
						r.Delete("/{trackId}", app.pageController.RemoveTrackFromPageHandler)
					})
					r.Route("/submissions", func(r chi.Router) {
						r.Get("/", app.pageController.GetSubmissionsHandler)
						r.Post("/", app.pageController.CreateSubmissionHandler)
						r.Post("/{submissionId}/approve", app.pageController.ApproveSubmissionHandler)
					})
				})
				r.Route("/feed", func(r chi.Router) {
					r.Get("/", app.pageController.GetPageFeedHandler)
				})
			})
		})

		r.Get("/artists/{artistId}/pages", app.pageController.GetArtistPagesHandler)
		r.Get("/tracks/pages", app.pageController.GetTrackPagesBulkHandler)
		r.Get("/tracks/uses", app.trackController.TrackUsesHandlerGet)
		r.Get("/tracks/{trackId}/pages", app.pageController.GetTrackPagesHandler)

		r.Route("/feed", func(r chi.Router) {
			r.Use(AuthJWTMiddleware)
			r.Get("/following", app.pageController.GetFollowingFeedHandler)
		})

		r.Route("/webhooks", func(r chi.Router) {
			r.Route("/stripe", func(r chi.Router) {
				r.Options("/", app.webhookController.StripeWebhookHandlerOptions)
				r.Post("/", app.webhookController.StripeWebhookHandlerPost)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", server.Addr)

	return server.ListenAndServe()
}
