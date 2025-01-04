package repository

import (
	"github.com/trentjkelly/layerr/internals/config"
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PortraitStorageRepository struct {
	r2Config		*aws.Config
	r2Client		*s3.Client
	portraitBucketName *string
}

// Constructor for new PortraitRepository
func NewPortraitStorageRepository() *PortraitStorageRepository {
	portraitStorageRepository := new(PortraitStorageRepository)
	portraitStorageRepository.r2Config = config.CreateR2Config()
	portraitStorageRepository.r2Client = config.CreateR2Client(portraitStorageRepository.r2Config)
	portraitStorageRepository.portraitBucketName = aws.String("artist-portrait")
	return portraitStorageRepository
}

// Uploads a portrait to storage
func (r *PortraitStorageRepository) CreatePortrait(ctx context.Context, file multipart.File, filename *string) error {
	
	input := &s3.PutObjectInput{
		Bucket:	r.portraitBucketName,
		Key:	filename,
		Body:	file,
	}

	_, err := r.r2Client.PutObject(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

// Gets a portrait from storage (to be streamed)
// func (r *PortraitStorageRepository) ReadPortrait() error {}

// Updates the portrait in storage
// func (r *PortraitStorageRepository) UpdatePortrait() error {}

// Deletes the portrait from storage
// func (r *PortraitStorageRepository) DeletePortrait() error {}
