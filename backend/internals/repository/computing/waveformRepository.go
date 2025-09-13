package computingRepository

import (
	"fmt"
	"math"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
)

type WaveformHeightsRepository struct {}

func NewWaveformHeightsRepository() *WaveformHeightsRepository {
	waveformHeightsRepository := new(WaveformHeightsRepository)
	return waveformHeightsRepository
}

func (r *WaveformHeightsRepository) CreateWaveform(flacPath string) ([]int, error) {
	file, err := os.Open(flacPath)
	if err != nil {
		return nil, fmt.Errorf("could not open the audio file: %w", err)
	}
	defer file.Close()

	streamer, _, err := flac.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("could not decode the flac file: %w", err)
	}
	defer streamer.Close()

	barsSlice, err := r.getBars(streamer)
	if err != nil {
		return nil, fmt.Errorf("could not get the bars from the song: %w", err)
	}
	
	return barsSlice, nil
}

func (r *WaveformHeightsRepository) getBars(streamer beep.StreamSeekCloser) ([]int, error) {
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

func (r *WaveformHeightsRepository) getSampleAmplitude(sample [2]float64) float64 {
	left, right := sample[0], sample[1]
	leftPercent := math.Abs(left) * 100
	rightPercent := math.Abs(right) * 100

	return math.Max(leftPercent, rightPercent)
}

func (r *WaveformHeightsRepository) getMaxFromSlice(slice []float64) float64 {
	max := slice[0]
	for _, v := range slice {
		if v > max {
			max = v
		}
	}
	return max
}

func (r *WaveformHeightsRepository) floatsToInts(floats []float64) []int {
	ints := make([]int, len(floats))
	for i, f := range floats {
		ints[i] = int(math.Round(f))
	}
	return ints
}

func (r *WaveformHeightsRepository) cleanBars(bars *[]int) {
	for i, value := range *bars {
		if value < 1 {
			(*bars)[i] = 1
		} else if value > 99 {
			(*bars)[i] = 100
		}
	}
}

