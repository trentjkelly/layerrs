package main

import (
	"log"
	"net/http"
	"time"
	"fmt"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// AWS SDK for Cloudflare R2 storage
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type appConfig struct {
	addr string
}

type application struct {
	config 		appConfig
	r2Config 	aws.Config
	r2Client 	*s3.Client
}

func createR2Config(accessKeyId string, accessKeySecret string) aws.Config {
	r2Config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),	
	)

	if err != nil {
		log.Fatal(err)
	}

	return r2Config
}

func createR2Client(r2Config aws.Config, accountId string) *s3.Client {
	r2Client := s3.NewFromConfig(r2Config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return r2Client
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	
	// Using middleware
	r.Use(middleware.RequestID)
  	r.Use(middleware.RealIP)
  	r.Use(middleware.Logger)
  	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 60))

	r.Get("/track", app.trackHandlerGet)
	r.Put("/track", app.trackHandlerPut)
	r.Options("/track", app.trackHandlerOptions)
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