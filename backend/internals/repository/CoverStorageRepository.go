package repository

import (
	"github.com/trentjkelly/layerrs/internals/config"
	"context"
	"mime/multipart"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type CoverStorageRepository struct {
	r2Config		*aws.Config
	r2Client		*s3.Client
	coverBucketName *string
}

// Constructor for new CoverRepository
func NewCoverStorageRepository() *CoverStorageRepository {
	coverStorageRepository := new(CoverStorageRepository)
	coverStorageRepository.r2Config = config.CreateR2Config()
	coverStorageRepository.r2Client = config.CreateR2Client(coverStorageRepository.r2Config)
	coverStorageRepository.coverBucketName = aws.String(os.Getenv("TRACK_COVER_BUCKET_NAME"))
	return coverStorageRepository
}

// Uploads a cover to storage
func (r *CoverStorageRepository) CreateCover(ctx context.Context, file multipart.File, filename *string) error {
	
	input := &s3.PutObjectInput{
		Bucket:	r.coverBucketName,
		Key:	filename,
		Body:	file,
	}

	_, err := r.r2Client.PutObject(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func (r *CoverStorageRepository) ReadCover(ctx context.Context, r2CoverKey *string) (io.ReadCloser, error) {

	input := &s3.GetObjectInput{
		Bucket: r.coverBucketName,
		Key: r2CoverKey,
	}

	res, err := r.r2Client.GetObject(ctx, input)

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// Gets a cover from storage (to be loaded to frontend)
// func (r *CoverStorageRepository) ReadCover() error {}

// Updates the cover in storage
// func (r *CoverStorageRepository) UpdateCover() error {}

// Deletes the cover from storage
// func (r *CoverStorageRepository) DeleteCover() error {}
