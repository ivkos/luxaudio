package analyzers

import (
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/fourier"
	"luxaudio/utils"
	"math"
	"math/cmplx"
)

type SmartAnalyzer struct {
	fftSize    int
	ledCount   int
	sampleRate float64

	intensities []float64

	freqs []float64
	loF   int
	hiF   int

	decayFactor   float64
	dbfsThreshold float64

	window []float64
	fft    *fourier.FFT

	mirror bool
}

func NewSmartAnalyzer(
	fftSize int,
	ledCount int,
	sampleRate float64,
	decayFactor float64,
	dbfsThreshold float64,
	audibleLow float64,
	audibleHigh float64,
	mirror bool,
) Analyzer {
	intensitiesLength := fftSize/2 + 1

	freqs := calculateFreqs(intensitiesLength, sampleRate, fftSize)

	return &SmartAnalyzer{
		fftSize:    fftSize,
		ledCount:   ledCount,
		sampleRate: sampleRate,

		intensities: make([]float64, intensitiesLength),

		freqs: freqs,
		loF:   getLowFreqIndex(freqs, audibleLow),
		hiF:   getHighFreqIndex(freqs, audibleHigh),

		decayFactor:   decayFactor,
		dbfsThreshold: dbfsThreshold,

		window: getHannWindow(fftSize),
		fft:    fourier.NewFFT(fftSize),

		mirror: mirror,
	}
}

func (sa *SmartAnalyzer) Analyze(sampleChunk []float64) []float64 {
	floats.Mul(sampleChunk, sa.window)
	ffs := sa.fft.Coefficients(nil, sampleChunk)

	for i := range sa.intensities {
		x := ffs[i]
		magnitude := cmplx.Abs(x)

		db := 20 * math.Log10(magnitude/(float64(sa.fftSize)/4))
		newIntensity := math.Min((math.Max(sa.dbfsThreshold, db)-sa.dbfsThreshold)/-sa.dbfsThreshold, 1)

		if sa.decayFactor != float64(0) && newIntensity <= sa.intensities[i] {
			sa.intensities[i] *= sa.decayFactor
		} else {
			sa.intensities[i] = newIntensity
		}
	}

	result := sa.intensities[sa.loF : sa.hiF+1]
	if sa.mirror {
		result = mirrorResult(result)
	}

	result = utils.ChunkedMean(result, sa.ledCount)
	result = utils.CenterArray(result, sa.ledCount)

	return result
}

func mirrorResult(original []float64) []float64 {
	reversed := append([]float64{}, original...)
	floats.Reverse(reversed)

	result := append([]float64{}, reversed...)
	result = append(result, original...)

	return result
}

func getLowFreqIndex(frequencies []float64, audibleLow float64) int {
	for i, f := range frequencies {
		if f >= audibleLow {
			return i
		}
	}

	return 0
}

func getHighFreqIndex(frequencies []float64, audibleHigh float64) int {
	for i, f := range frequencies {
		if f >= audibleHigh {
			return i
		}
	}

	return len(frequencies) - 1
}

func calculateFreqs(size int, sampleRate float64, fftSize int) []float64 {
	freqs := make([]float64, size)
	coef := sampleRate / float64(fftSize)
	for i := range freqs {
		freqs[i] = float64(i) * coef
	}
	return freqs
}

func getHannWindow(size int) []float64 {
	r := make([]float64, size)

	if size == 1 {
		r[0] = 1
	} else {
		N := size - 1
		coef := 2 * math.Pi / float64(N)
		for n := 0; n <= N; n++ {
			r[n] = 0.5 * (1 - math.Cos(coef*float64(n)))
		}
	}

	return r
}
