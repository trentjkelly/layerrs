package computingRepository

import (
	"fmt"
	"math"
	"context"
	"mime/multipart"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

type WaveformComputingRepository struct {}

func NewWaveformComputingRepository() *WaveformComputingRepository {
	waveformComputingRepository := new(WaveformComputingRepository)
	return waveformComputingRepository
}

func (r *WaveformComputingRepository) CreateWaveform(ctx context.Context, audio multipart.File) ([]int, error) {
	streamer, _, err := mp3.Decode(audio)
	if err != nil {
		return nil, fmt.Errorf("could not decode the mp3 file: %w", err)
	}
	defer streamer.Close()

	barsSlice, err := r.getBars(streamer)
	if err != nil {
		return nil, fmt.Errorf("could not get the bars from the song: %w", err)
	}
	
	return barsSlice, nil
}

func (r *WaveformComputingRepository) getBars(streamer beep.StreamSeekCloser) ([]int, error) {
	numSamples := streamer.Len()
	sampleSlice := make([][2]float64, numSamples)

	n, ok := streamer.Stream(sampleSlice)
	if !ok || n == 0 {
		return nil, fmt.Errorf("could not get the amplitudes from the sample slice")
	}

	amplitudeSlice := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		amplitudeSlice[i] = r.getSampleAmplitude(sampleSlice[i])
	}

	numBars := 300
	chunkSize := numSamples / numBars
	barsAmplitudeSlice := make([]float64, numBars)

	for i := 0; i < numBars; i++ {
		lowerBound := i * chunkSize
		upperBound := (i + 1) * chunkSize
		barsAmplitudeSlice[i] = r.getMaxFromSlice(amplitudeSlice[lowerBound:upperBound])
	}

	barsSlice := r.floatsToInts(barsAmplitudeSlice)
	r.cleanBars(&barsSlice)
	return barsSlice, nil
}

func (r *WaveformComputingRepository) getSampleAmplitude(sample [2]float64) float64 {
	left, right := sample[0], sample[1]
	leftPercent := math.Abs(left) * 100
	rightPercent := math.Abs(right) * 100

	return math.Max(leftPercent, rightPercent)
}

func (r *WaveformComputingRepository) getMaxFromSlice(slice []float64) float64 {
	max := slice[0]
	for _, v := range slice {
		if v > max {
			max = v
		}
	}
	return max
}

func (r *WaveformComputingRepository) floatsToInts(floats []float64) []int {
	ints := make([]int, len(floats))
	for i, f := range floats {
		ints[i] = int(math.Round(f))
	}
	return ints
}

func (r *WaveformComputingRepository) cleanBars(bars *[]int) {
	for i, value := range *bars {
		if value < 1 {
			(*bars)[i] = 1
		} else if value > 99 {
			(*bars)[i] = 100
		}
	}
}

