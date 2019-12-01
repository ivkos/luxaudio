package utils

import (
	"flag"
	"os"
)

func GetFlags() (string, uint16, int, int, int, int, float64) {
	var host = flag.String("host", "", "host of the luxsrv")
	var port = flag.Uint("port", DefaultPort, "port of the luxsrv")
	var ledCount = flag.Int("leds", 0, "number of LEDs to be driven (max 255)")
	var fftSize = flag.Int("fft", 2048, "FFT size")
	var sampleRate = flag.Int("sampleRate", 0, "sample rate")
	var channels = flag.Int("channels", 2, "number of channels")
	var decay = flag.Float64("decay", 0.50, "decay factor [0,1] controls the smoothness of the visualization")

	flag.Parse()

	if *host == "" || *ledCount == 0 || *ledCount > 255 || *sampleRate == 0 {
		flag.Usage()
		os.Exit(2)
	}

	return *host, uint16(*port), *ledCount, *fftSize, *sampleRate, *channels, *decay
}
