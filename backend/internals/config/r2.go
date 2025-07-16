package config

import (
	"log"
	"context"
	"os"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

)

func CreateR2Config() *aws.Config {
	accessKeyId := os.Getenv("R2_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("R2_SECRET_ACCESS_KEY_ID")

	r2Config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),	
	)
	if err != nil {
		log.Fatal(err)
	}

	return &r2Config
}

func CreateR2Client(r2Config *aws.Config) *s3.Client {
	accountId := os.Getenv("R2_ACCOUNT_ID")

	r2Client := s3.NewFromConfig(*r2Config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return r2Client
}
