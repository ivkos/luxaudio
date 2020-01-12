package analyzers

import (
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
	"luxaudio/utils"
	"math"
)

type SmartAnalyzer struct {
	fftSize    int
	ledCount   int
	sampleRate float64

	intensities []float64
	ledData     []byte

	freqsInit bool
	loF       int
	hiF       int

	freqs []float64

	decayFactor   float64
	dbfsThreshold float64
}

func NewSmartAnalyzer(fftSize int, ledCount int, sampleRate float64, decayFactor float64, dbfsThreshold float64) Analyzer {
	return &SmartAnalyzer{
		fftSize:    fftSize,
		ledCount:   ledCount,
		sampleRate: sampleRate,

		intensities: make([]float64, fftSize/2),
		ledData:     make([]byte, ledCount*3),

		freqsInit: false,

		decayFactor:   decayFactor,
		dbfsThreshold: dbfsThreshold,
	}
}

func (sa *SmartAnalyzer) Analyze(sampleChunk []float64) []byte {
	window.Apply(sampleChunk, window.Hann)
	ffs := fft.FFTReal(sampleChunk)

	for i := range sa.intensities {
		x := ffs[i]
		magnitude := real(x)*real(x) + imag(x)*imag(x)

		db := 10 * math.Log10(magnitude/math.Pow(float64(sa.fftSize)/2, 2))
		newIntensity := math.Min((math.Max(-sa.dbfsThreshold, db)+sa.dbfsThreshold)/sa.dbfsThreshold, 1)

		if sa.decayFactor != float64(0) && newIntensity <= sa.intensities[i] {
			sa.intensities[i] *= sa.decayFactor
		} else {
			sa.intensities[i] = newIntensity
		}
	}

	if !sa.freqsInit {
		sa.freqs = sa.calculateFreqs()
		sa.loF = getLowFreqIndex(sa.freqs)
		sa.hiF = getHighFreqIndex(sa.freqs)
		sa.freqsInit = true
	}

	spectrum := utils.ChunkedMean(sa.intensities[sa.loF:sa.hiF], sa.ledCount)
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

func (sa *SmartAnalyzer) calculateFreqs() []float64 {
	freqs := make([]float64, sa.fftSize/2)
	coef := sa.sampleRate / float64(sa.fftSize)
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
