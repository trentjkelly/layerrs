package main

import (
	"log"

	"github.com/stripe/stripe-go/v85"
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
	PRODUCTION  = "PRODUCTION"
)

const (
	WEBP_IMAGE_QUALITY = 85
)

func main() {
	// Get the environment
	log.Println("Building the application")
	env, isDocker, err := config.GetEnvironment()
	if err != nil {
		log.Fatalf("Could not get the environment: %v", err)
	}

	// Database Connection
	log.Println("Connecting to the database")
	pool, err := config.InitDB(env, isDocker)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	var frontendUrl string
	var magicLink string
	if env == DEVELOPMENT {
		magicLink = "http://localhost:8080/api/authentication/verify?token="
		frontendUrl = "https://localhost:3000"
	} else {
		magicLink = "https://layerrs.com/api/authentication/verify?token="
		frontendUrl = "https://layerrs.com"
	}

	// Setting global Stripe API Key
	stripeConfig, err := config.NewStripeConfig(env, frontendUrl)
	if err != nil {
		log.Fatalf("Could not create the Stripe config: %v", err)
	}
	stripe.Key = stripeConfig.SecretKey

	// -- REPOSITORIES --
	log.Println("Creating the repositories, services, and controllers")

	// Computing Repositories
	trackConversionRepo := computingRepository.NewTrackConversionRepository()
	waveformRepo := computingRepository.NewWaveformHeightsRepository()
	portraitConversionRepo := computingRepository.NewPortraitConversionRepository(WEBP_IMAGE_QUALITY)

	// Auth Repositories
	passwordRepo := authRepository.NewPasswordRepository()
	authRepo := authRepository.NewAuthRepository()
	verificationEmailRepo := authRepository.NewVerificationEmailRepository()

	// Database Repositories
	artistDatabaseRepo := databaseRepository.NewArtistDatabaseRepository(pool)
	likesDatabaseRepo := databaseRepository.NewLikesDatabaseRepository(pool)
	trackDatabaseRepo := databaseRepository.NewTrackDatabaseRepository(pool)
	trackTreeDatabaseRepo := databaseRepository.NewTrackTreeDatabaseRepository(pool)
	waveformDatabaseRepo := databaseRepository.NewWaveformDatabaseRepository(pool)
	layerrsDatabaseRepo := databaseRepository.NewLayerrsDatabaseRepository(pool)
	authDatabaseRepo := databaseRepository.NewAuthDatabaseRepository(pool)
	purchaseDatabaseRepo := databaseRepository.NewPurchaseRepository(pool)

	// Storage Repositories
	portraitStorageRepo := storageRepository.NewPortraitStorageRepository(env)
	trackStorageRepo := storageRepository.NewTrackStorageRepository(env)

	// -- SERVICES --
	authService := service.NewAuthService(passwordRepo, artistDatabaseRepo, authRepo, verificationEmailRepo, authDatabaseRepo, magicLink)
	trackService := service.NewTrackService(trackStorageRepo, portraitStorageRepo, trackDatabaseRepo, trackTreeDatabaseRepo, trackConversionRepo, waveformRepo, waveformDatabaseRepo, layerrsDatabaseRepo, env)
	recService := service.NewRecommendationsService(trackDatabaseRepo, likesDatabaseRepo, portraitStorageRepo)
	artistService := service.NewArtistService(artistDatabaseRepo, portraitStorageRepo, portraitConversionRepo)
	likesService := service.NewLikesService(likesDatabaseRepo, trackDatabaseRepo)
	layerrsService := service.NewLayerrsService(layerrsDatabaseRepo, portraitStorageRepo)
	sellerService := service.NewSellerService(artistDatabaseRepo, stripeConfig.ReturnURL, stripeConfig.RefreshURL)
	purchaseService := service.NewPurchaseService(purchaseDatabaseRepo, trackStorageRepo, stripeConfig.ReturnURL, stripeConfig.RefreshURL)

	// -- CONTROLLERS --
	authController := controller.NewAuthController(authService, frontendUrl)
	trackController := controller.NewTrackController(trackService)
	recController := controller.NewRecommendationsController(recService)
	likesController := controller.NewLikesController(likesService)
	artistController := controller.NewArtistController(artistService)
	layerrsController := controller.NewLayerrsController(layerrsService)
	sellerController := controller.NewSellerController(sellerService)
	webhookController := controller.NewWebhookController(sellerService, stripeConfig.WebhookSecret)
	purchaseController := controller.NewPurchaseController(purchaseService)

	// -- CONFIGURATION --
	cfg := appConfig{
		addr: ":8080",
	}

	app := &application{
		config:                    cfg,
		trackController:           trackController,
		recommendationsController: recController,
		authController:            authController,
		likesController:           likesController,
		artistController:          artistController,
		layerrsController:         layerrsController,
		sellerController:          sellerController,
		webhookController:         webhookController,
		purchaseController:        purchaseController,
	}

	// Mount and run the application
	log.Println("Starting the application")
	mux := app.mount()
	log.Fatal(app.run(mux))
}
