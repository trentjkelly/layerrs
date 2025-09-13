package computingRepository

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"os/exec"
	"mime/multipart"
	"os"
	"path/filepath"
)

type TrackConversionRepository struct {}

func NewTrackConversionRepository() *TrackConversionRepository {
	trackConversionRepository := new(TrackConversionRepository)
	return trackConversionRepository
}

func (r *TrackConversionRepository) ConvertAllTracks(audio multipart.File, trackId string, audioExtension string) (string, string, string, string, string, string, error) {
	tempDirPath, err := os.MkdirTemp("", trackId)
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("failed to create temp path: %w", err)
	}

	tempFilePath := filepath.Join(tempDirPath, fmt.Sprintf("%s.%s", trackId, audioExtension))
	if tempFilePath == "" {
		return "", "", "", "", "", "", fmt.Errorf("tempFilePath could not be created")
	}

	err = r.WriteFileToTempPath(audio, tempFilePath)
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("could not write file: %w", err)
	}

	flacPath, aacPath, opusPath, flacName, opusName, aacName := r.CreatePathNames(tempDirPath, trackId)

	err = r.FFMPEGConversions(tempFilePath, flacPath, opusPath, aacPath) 
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("could not convert tracks to different types: %w", err)
	}

	return flacPath, opusPath, aacPath, flacName, opusName, aacName, nil
}

func (r *TrackConversionRepository) CreatePathNames(tempDirPath string, trackId string) (string, string, string, string, string, string) {

	flacName := fmt.Sprintf("%s.flac", trackId)
	opusName := fmt.Sprintf("%s.opus", trackId)
	aacName := fmt.Sprintf("%s.aac", trackId)
	flacPath := filepath.Join(tempDirPath, flacName)
	aacPath := filepath.Join(tempDirPath, aacName)
	opusPath := filepath.Join(tempDirPath, opusName)
	return flacPath, aacPath, opusPath, flacName, opusName, aacName
}

func (r *TrackConversionRepository) FFMPEGConversions(inputPath string, flacPath string, opusPath string, aacPath string) error {
	err := r.ConvertTrackToFLAC(inputPath, flacPath)
	if err != nil {
		return fmt.Errorf("failed to convert track to FLAC: %w", err)
	}
	
	err = r.ConvertTrackToAAC(flacPath, aacPath)
	if err != nil {
		return fmt.Errorf("failed to convert track to AAC: %w", err)
	}
	
	
	err = r.ConvertTrackToOPUS(flacPath, opusPath)
	if err != nil {
		return fmt.Errorf("failed to convert track to OPUS: %w", err)
	}

	err = os.Remove(inputPath)
	if err != nil {
		return fmt.Errorf("failed to remove input file: %w", err)
	}

	return nil
}

func (r *TrackConversionRepository) WriteFileToTempPath(audio multipart.File, tempPath string) error {
	destinationTempFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("could not create temp file for audio: %w", err)
	}
	defer destinationTempFile.Close()

	_, err = io.Copy(destinationTempFile, audio)
	return err
}

func (r *TrackConversionRepository) GetAACTrackDuration(aacPath string) (float64, error) {
	cmd := exec.Command(
		"ffprobe", 
		"-i", 
		aacPath, 
		"-show_entries", 
		"format=duration", 
		"-v", 
		"quiet", 
		"-of", 
		"csv=p=0",
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get the duration of the AAC track: %w", err)
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse the duration of the AAC track: %w", err)
	}

	return duration, nil
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
