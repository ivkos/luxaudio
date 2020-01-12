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
	ledData     []byte

	freqs []float64
	loF   int
	hiF   int

	decayFactor   float64
	dbfsThreshold float64

	window []float64
	fft    *fourier.FFT
}

func NewSmartAnalyzer(fftSize int, ledCount int, sampleRate float64, decayFactor float64, dbfsThreshold float64) Analyzer {
	intensitiesLength := fftSize/2 + 1

	freqs := calculateFreqs(intensitiesLength, sampleRate, fftSize)

	return &SmartAnalyzer{
		fftSize:    fftSize,
		ledCount:   ledCount,
		sampleRate: sampleRate,

		intensities: make([]float64, intensitiesLength),
		ledData:     make([]byte, ledCount*3),

		freqs: freqs,
		loF:   getLowFreqIndex(freqs),
		hiF:   getHighFreqIndex(freqs),

		decayFactor:   decayFactor,
		dbfsThreshold: dbfsThreshold,

		window: getHannWindow(fftSize),
		fft:    fourier.NewFFT(fftSize),
	}
}

func (sa *SmartAnalyzer) Analyze(sampleChunk []float64) []byte {
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

	spectrum := utils.ChunkedMean(sa.intensities[sa.loF:sa.hiF+1], sa.ledCount)
	spectrum = utils.CenterArray(spectrum, sa.ledCount)

	var r, g, b float64 = 255, 0, 255
	for i, x := range spectrum {
		sa.ledData[i*3+0] = byte(g * x)
		sa.ledData[i*3+1] = byte(r * x)
		sa.ledData[i*3+2] = byte(b * x)
	}

	return sa.ledData
}

const hearingRangeLow = 20
const hearingRangeHigh = 20000

func getLowFreqIndex(frequencies []float64) int {
	for i, f := range frequencies {
		if f >= hearingRangeLow {
			return i
		}
	}

	return 0
}

func getHighFreqIndex(frequencies []float64) int {
	for i, f := range frequencies {
		if f >= hearingRangeHigh {
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

func Ra(f float64) float64 {
	return (math.Pow(12194, 2) * math.Pow(f, 4)) /
		((math.Pow(f, 2) + math.Pow(20.6, 2)) *
			math.Sqrt((math.Pow(f, 2)+math.Pow(107.7, 2))*
				(math.Pow(f, 2)+math.Pow(737.9, 2))) *
			(math.Pow(f, 2) + math.Pow(12194, 2)))
}

func Aw(f float64) float64 {
	return 20*math.Log10(Ra(f)) - 20*math.Log10(Ra(1000))
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
