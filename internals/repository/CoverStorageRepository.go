package repository

import (
	"github.com/trentjkelly/layerr/internals/config"
	"context"
	"log"
	"mime/multipart"

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
	coverStorageRepository.coverBucketName = aws.String("cover-art")
	return coverStorageRepository
}

// Uploads a cover to storage
func (r *CoverStorageRepository) CreateCover(file multipart.File, filename *string) error {
	
	input := &s3.PutObjectInput{
		Bucket:	r.coverBucketName,
		Key:	filename,
		Body:	file,
	}

	res, err := r.r2Client.PutObject(context.TODO(), input)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("File uploaded!")
	log.Println(res)

	return nil
}

// Gets a cover from storage (to be loaded to frontend)
// func (r *CoverStorageRepository) ReadCover() error {}

// Updates the cover in storage
// func (r *CoverStorageRepository) UpdateCover() error {}

// Deletes the cover from storage
// func (r *CoverStorageRepository) DeleteCover() error {}
