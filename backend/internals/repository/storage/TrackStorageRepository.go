package storageRepository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"log"
	"strings"
	"time"

	"github.com/trentjkelly/layerrs/internals/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type TrackStorageRepository struct {
	r2Config		*aws.Config
	r2Client		*s3.Client
	r2Presigner		*s3.PresignClient
	trackWavBucketName *string
	trackAacBucketName  *string
	environment		string
}

// Constructor for new TrackRepository
func NewTrackStorageRepository(environment string) *TrackStorageRepository {
	trackStorageRepository := new(TrackStorageRepository)
	trackStorageRepository.r2Config = config.CreateR2Config()
	trackStorageRepository.r2Client = config.CreateR2Client(trackStorageRepository.r2Config)
	trackStorageRepository.r2Presigner = config.CreateR2Presigner(trackStorageRepository.r2Client)
	trackStorageRepository.trackWavBucketName = aws.String(os.Getenv(fmt.Sprintf("TRACK_AUDIO_WAV_BUCKET_NAME_%s", environment)))
	trackStorageRepository.trackAacBucketName = aws.String(os.Getenv(fmt.Sprintf("TRACK_AUDIO_AAC_BUCKET_NAME_%s", environment)))
	trackStorageRepository.environment = environment
	return trackStorageRepository
}

// Uploads all tracks to R2
func (r *TrackStorageRepository) CreateAllTracks(ctx context.Context, wavPath string, aacPath string, wavKey string, aacKey string) error {
	envName := strings.ToLower(r.environment)
	
	// Upload WAV file
	file, err := os.Open(wavPath)
	if err != nil {
		return fmt.Errorf("could not open the wav file: %w", err)
	}

	wavBucketName := fmt.Sprintf("track-audio-wav-%s", envName)
	err = r.CreateTrack(ctx, file, wavKey, wavBucketName)
	if err != nil {
		return fmt.Errorf("failed to upload wav file to R2: %w", err)
	}
	file.Close()

	// Upload AAC file
	file, err = os.Open(aacPath)
	if err != nil {
		return fmt.Errorf("could not open the aac file: %w", err)
	}

	aacBucketName := fmt.Sprintf("track-audio-aac-%s", envName)
	err = r.CreateTrack(ctx, file, aacKey, aacBucketName)
	if err != nil {
		return fmt.Errorf("failed to upload wav file to R2: %w", err)
	}
	file.Close()

	return nil
}

// Uploads a single track to R2
func (r *TrackStorageRepository) CreateTrack(ctx context.Context, file multipart.File, filename string, bucketName string) error {
	input := &s3.PutObjectInput{
		Bucket:	&bucketName,
		Key:	&filename,
		Body:	file,
	}

	_, err := r.r2Client.PutObject(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

// Gets a track from storage (to be streamed)
func (r *TrackStorageRepository) ReadWavTrack(ctx context.Context, trackName *string, startByte int, endByte int) (io.ReadCloser, error) {
	rangeString := fmt.Sprintf("bytes=%d-%d", startByte, endByte)

	input := &s3.GetObjectInput{
		Bucket: r.trackWavBucketName,
		Key: trackName,
		Range: aws.String(rangeString),
	}

	res, err := r.r2Client.GetObject(ctx, input)

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// Gets a signed url for a track
func (r *TrackStorageRepository) GetSignedAacURL(ctx context.Context, objectKey string, expirationTime time.Duration) (string, time.Duration, error) {

	log.Println("[DEBUG] Getting signed url for track: ", objectKey)
	input := &s3.GetObjectInput{
		Bucket: r.trackAacBucketName,
		Key: &objectKey,
	}

	req, err := r.r2Presigner.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expirationTime
	})
	if err != nil {
		return "", 0, fmt.Errorf("failed to get presigned url: %w", err)
	}

	log.Println("[DEBUG] Signed url for track: ", req.URL)

	return req.URL, expirationTime, nil
}

// Gets a signed url for a track
func ( r*TrackStorageRepository) GetSignedWavURL(ctx context.Context, objectKey string, expirationTime time.Duration) (string, time.Duration, error) {

	log.Println("[DEBUG] Getting signed url for track: ", objectKey)
	input := &s3.GetObjectInput{
		Bucket: r.trackWavBucketName,
		Key: &objectKey,
	}

	req, err := r.r2Presigner.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expirationTime
	})
	if err != nil {
		return "", 0, fmt.Errorf("failed to get presigned url: %w", err)
	}

	return req.URL, expirationTime, nil
}

// Updates the track in storage
// func (r *TrackStorageRepository) UpdateTrack() error {}

// Deletes the track from storage
// func (r *TrackStorageRepository) DeleteTrack() error {}
