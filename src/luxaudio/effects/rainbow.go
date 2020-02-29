package effects

import (
	"image/color"
	"time"
)

type RainbowEffect struct {
	rate     float64
	ledCount int
	ledData  []byte
	rainbow  []color.RGBA
}

func NewRainbowEffect(ledCount int, rate float64) Effect {
	e := &RainbowEffect{
		rate:     rate,
		ledCount: ledCount,
		ledData:  make([]byte, ledCount*3),
		rainbow:  make([]color.RGBA, ledCount*3),
	}

	go e.startRainbow()

	return e
}

func (e *RainbowEffect) Apply(intensities []float64) []byte {
	for i, x := range intensities {
		e.ledData[i*3+0] = byte(float64(e.rainbow[i].G) * x)
		e.ledData[i*3+1] = byte(float64(e.rainbow[i].R) * x)
		e.ledData[i*3+2] = byte(float64(e.rainbow[i].B) * x)
	}

	return e.ledData
}

func (e *RainbowEffect) startRainbow() {
	offset := 0

	for {
		t := time.NewTimer(time.Duration(1000/e.rate) * time.Millisecond)
		<-t.C

		for i := 0; i < e.ledCount; i++ {
			e.rainbow[i] = wheel(uint8((offset + i) & 255))
		}

		offset = (offset + 1) % 256
	}
}

func wheel(pos uint8) color.RGBA {
	if pos < 85 {
		return color.RGBA{
			R: pos * 3,
			G: 255 - pos*3,
			B: 0,
		}
	}

	if pos < 170 {
		pos -= 85
		return color.RGBA{
			R: 255 - pos*3,
			G: 0,
			B: pos * 3,
		}
	}

	pos -= 170
	return color.RGBA{
		R: 0,
		G: pos * 3,
		B: 255 - pos*3,
	}
}
