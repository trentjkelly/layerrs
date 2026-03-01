package storageRepository

import (
	"github.com/trentjkelly/layerrs/internals/config"
	"context"
	"fmt"
	"bytes"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PortraitStorageRepository struct {
	r2Config		*aws.Config
	r2Client		*s3.Client
	r2Presigner		*s3.PresignClient
	portraitBucketName *string
}

// Constructor for new PortraitRepository
func NewPortraitStorageRepository(env string) *PortraitStorageRepository {
	portraitStorageRepository := new(PortraitStorageRepository)
	portraitStorageRepository.r2Config = config.CreateR2Config()
	portraitStorageRepository.r2Client = config.CreateR2Client(portraitStorageRepository.r2Config)
	portraitStorageRepository.r2Presigner = config.CreateR2Presigner(portraitStorageRepository.r2Client)
	portraitStorageRepository.portraitBucketName = aws.String(os.Getenv(fmt.Sprintf("ARTIST_PORTRAIT_BUCKET_NAME_%s", env)))
	return portraitStorageRepository
}

// Uploads a portrait to storage
func (r *PortraitStorageRepository) CreatePortrait(ctx context.Context, file *bytes.Reader, filename string) error {
	input := &s3.PutObjectInput{
		Bucket:	r.portraitBucketName,
		Key:	&filename,
		Body:	file,
	}

	_, err := r.r2Client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// Gets a pre-signed URL for a portrait
func (r *PortraitStorageRepository) GetSignedPortraitURL(ctx context.Context, key string, expirationTime time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: r.portraitBucketName,
		Key:    &key,
	}

	req, err := r.r2Presigner.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expirationTime
	})
	if err != nil {
		return "", fmt.Errorf("failed to get presigned portrait url: %w", err)
	}

	return req.URL, nil
}

// Gets a portrait from storage (to be streamed)
// func (r *PortraitStorageRepository) ReadPortrait() error {}

// Updates the portrait in storage
// func (r *PortraitStorageRepository) UpdatePortrait() error {}

// Deletes the portrait from storage
// func (r *PortraitStorageRepository) DeletePortrait() error {}
