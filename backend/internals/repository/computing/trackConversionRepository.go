package computingRepository

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

type TrackConversionRepository struct {}

func NewTrackConversionRepository() *TrackConversionRepository {
	trackConversionRepository := new(TrackConversionRepository)
	return trackConversionRepository
}

// Converts the original track to all formats, returns the filepaths to upload to R2
func (r *TrackConversionRepository) ConvertAllFormats(audio multipart.File, audioFileExtension string, artistId string, trackId string) (string, string, string, string, string, string, error) {
	// Create temp directory of os.Temp() + artist id
	tempPath := filepath.Join(os.TempDir(), artistId)
	err := os.MkdirAll(tempPath, 0755)
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Write normal file to the temp directory
	originalFileName := fmt.Sprintf("%s%s", "original", audioFileExtension)
	originalFilePath := filepath.Join(tempPath, originalFileName)

	tempAudioFile, err := os.Create(originalFilePath)
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("failed to create file for the original track: %w", err)
	}
	defer tempAudioFile.Close()

	_, err = io.Copy(tempAudioFile, audio)
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("failed to copy audio file into the original path: %w", err)
	}

	// Apply conversion of original to FLAC (new track name + .flac)
	flacFileName := fmt.Sprintf(trackId, ".flac")
	flacFilePath := filepath.Join(tempPath, flacFileName)

	err = r.ConvertTrackToFLAC(originalFilePath, flacFilePath)
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("failed to convert original track to flac: %w", err)
	}

	// Apply conversion of FLAC to OPUS
	opusFileName := fmt.Sprintf(trackId, ".opus")
	opusFilePath := filepath.Join(tempPath, opusFileName)

	err = r.ConvertTrackToOPUS(originalFilePath, opusFilePath)
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("failed to convert original track to opus: %w", err)
	}

	// Apply conversion of FLAC to AAC
	aacFileName := fmt.Sprintf(trackId, ".aac")
	aacFilePath := filepath.Join(tempPath, aacFileName)

	err = r.ConvertTrackToAAC(originalFilePath, aacFilePath)
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("failed to convert original track to aac: %w", err)
	}

	// Return paths of FLAC, OPUS, and AAC
	return flacFilePath, opusFilePath, aacFilePath, flacFileName, opusFileName, aacFileName, nil
}

// Converts an audio file to AAC format
func (r *TrackConversionRepository) ConvertTrackToAAC(inputPath string, outputPath string) error {
	cmd := exec.Command(
		"ffmpeg", 
		"-i", 
		inputPath, 
		"-codec:a", 
		"libfdk_aac", 
		"-b:a", 
		"256k", 
		"-profile:a", 
		"aac_low", 
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert track to AAC: %w", err)
	}

	return nil
}

// Converts a track to FLAC format
func (r *TrackConversionRepository) ConvertTrackToFLAC(inputPath string, outputPath string) error {
	cmd := exec.Command(
		"ffmpeg", 
		"-i", 
		inputPath, 
		"-codec:a", 
		"flac", 
		"-compression_level", 
		"8",
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert track to FLAC: %w", err)
	}

	return nil
}

// Converts a track to OPUS format
func (r *TrackConversionRepository) ConvertTrackToOPUS(inputPath string, outputPath string) error {
	cmd := exec.Command(
		"ffmpeg", 
		"-i", 
		inputPath, 
		"-codec:a", 
		"libopus", 
		"-b:a", 
		"256k", 
		"-vbr", 
		"on", 
		"-application", 
		"audio", 
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert track to OPUS: %w", err)
	}

	return nil
}