package effects

import "image/color"

type SolidColorEffect struct {
	color color.RGBA

	ledCount int
	ledData  []byte
}

func NewSolidColorEffect(ledCount int, color color.RGBA) Effect {
	return &SolidColorEffect{
		color: color,

		ledCount: ledCount,
		ledData:  make([]byte, ledCount*3),
	}
}

func (e *SolidColorEffect) Apply(intensities []float64) []byte {
	for i, x := range intensities {
		e.ledData[i*3+0] = byte(float64(e.color.G) * x)
		e.ledData[i*3+1] = byte(float64(e.color.R) * x)
		e.ledData[i*3+2] = byte(float64(e.color.B) * x)
	}

	return e.ledData
}
