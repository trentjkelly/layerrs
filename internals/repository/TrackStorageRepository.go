package repository

import (
	"github.com/trentjkelly/layerr/internals/config"
	"context"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type TrackStorageRepository struct {
	r2Config		*aws.Config
	r2Client		*s3.Client
	trackBucketName *string
}

// Constructor for new TrackRepository
func NewTrackStorageRepository() *TrackStorageRepository {
	trackStorageRepository := new(TrackStorageRepository)
	trackStorageRepository.r2Config = config.CreateR2Config()
	trackStorageRepository.r2Client = config.CreateR2Client(trackStorageRepository.r2Config)
	trackStorageRepository.trackBucketName = aws.String("track-audio")
	return trackStorageRepository
}

// Uploads a track to storage
func (r *TrackStorageRepository) CreateTrack(file multipart.File, filename *string) error {
	
	input := &s3.PutObjectInput{
		Bucket:	r.trackBucketName,
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

// Gets a track from storage (to be streamed)
func (r *TrackStorageRepository) ReadTrack() error {

}

// Updates the track in storage
func (r *TrackStorageRepository) UpdateTrack() error {

}

// Deletes the track from storage
func (r *TrackStorageRepository) DeleteTrack() error {

}


