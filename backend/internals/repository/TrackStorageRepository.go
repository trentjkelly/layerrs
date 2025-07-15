package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/trentjkelly/layerrs/internals/config"

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
	trackStorageRepository.trackBucketName = aws.String(os.Getenv("TRACK_AUDIO_BUCKET_NAME"))
	return trackStorageRepository
}

// Uploads a track to storage
func (r *TrackStorageRepository) CreateTrack(ctx context.Context, file multipart.File, filename *string) error {
	
	input := &s3.PutObjectInput{
		Bucket:	r.trackBucketName,
		Key:	filename,
		Body:	file,
	}

	_, err := r.r2Client.PutObject(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

// Gets a track from storage (to be streamed)
func (r *TrackStorageRepository) ReadTrack(ctx context.Context, trackName *string, startByte int, endByte int) (io.ReadCloser, error) {
	rangeString := fmt.Sprintf("bytes=%d-%d", startByte, endByte)

	input := &s3.GetObjectInput{
		Bucket: r.trackBucketName,
		Key: trackName,
		Range: aws.String(rangeString),
	}

	res, err := r.r2Client.GetObject(ctx, input)

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// Updates the track in storage
// func (r *TrackStorageRepository) UpdateTrack() error {}

// Deletes the track from storage
// func (r *TrackStorageRepository) DeleteTrack() error {}
