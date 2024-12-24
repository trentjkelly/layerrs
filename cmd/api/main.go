package main

import (
	// Base imports
	"log"
	"os"
	
	// AWS SDK for Cloudflare R2 storage (might not need since abstracted in api.go)
	// "github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/credentials"
	// "github.com/aws/aws-sdk-go-v2/service/s3"

	// Loads environment variables from local .env file
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	// Cloudflare R2 Storage configuration
	r2Config := createR2Config(os.Getenv("R2_ACCESS_KEY_ID"), os.Getenv("R2_SECRET_ACCESS_ID"))
	r2Client := createR2Client(r2Config, os.Getenv("R2_ACCOUNT_ID"))

	// Setup configuration and injected dependencies
	cfg := appConfig{
		addr : ":8080",
	}

	app := &application{
		config		: cfg,
		r2Config	: r2Config, 
		r2Client	: r2Client,
	}

	// Mount and run the application
	mux := app.mount()
	log.Fatal(app.run(mux))
}
