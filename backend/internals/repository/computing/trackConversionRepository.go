package computingRepository

import (
	"fmt"
	"os/exec"
	"mime/multipart"
)

type TrackConversionRepository struct {}

func NewTrackConversionRepository() *TrackConversionRepository {
	trackConversionRepository := new(TrackConversionRepository)
	return trackConversionRepository
}

func (r *TrackConversionRepository) ConvertAllFormats(audio multipart.File, outputPath string) error {
	err := r.ConvertTrackToFLAC(inputPath, outputPath)
	if err != nil {
		return fmt.Errorf("failed to convert track to FLAC: %w", err)
	}
	
	err = r.ConvertTrackToAAC(inputPath, outputPath)
	if err != nil {
		return fmt.Errorf("failed to convert track to AAC: %w", err)
	}
	
	
	err = r.ConvertTrackToOPUS(inputPath, outputPath)
	if err != nil {
		return fmt.Errorf("failed to convert track to OPUS: %w", err)
	}

	return nil
}

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