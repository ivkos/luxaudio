package analyzers

import (
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/spectral"
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

	pwelchOptions *spectral.PwelchOptions
	freqs         []float64

	decayFactor float64
}

func NewSmartAnalyzer(fftSize int, ledCount int, sampleRate float64, decayFactor float64) Analyzer {
	return &SmartAnalyzer{
		fftSize:    fftSize,
		ledCount:   ledCount,
		sampleRate: sampleRate,

		intensities: make([]float64, fftSize/2),
		ledData:     make([]byte, ledCount*3),

		freqsInit: false,

		pwelchOptions: &spectral.PwelchOptions{
			NFFT:      fftSize,
			Window:    window.Hann,
			Scale_off: true,
		},

		decayFactor: decayFactor,
	}
}

func (sa *SmartAnalyzer) Analyze(sampleChunk []float64) []byte {
	window.Apply(sampleChunk, window.Hann)
	ffs := fft.FFTReal(sampleChunk)

	for i := range sa.intensities {
		x := ffs[i]
		magnitude := real(x)*real(x) + imag(x)*imag(x)

		db := 10 * math.Log10(magnitude/math.Pow(float64(sa.fftSize), 2))
		newIntensity := math.Min((math.Max(-75, db)+75)/75, 1)

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

	for i, x := range spectrum {
		sa.ledData[i*3+0] = byte(128 * x)
		sa.ledData[i*3+1] = byte(255 * x)
		sa.ledData[i*3+2] = byte(64 * x)
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
