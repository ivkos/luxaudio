package utils

import (
	"flag"
	"image/color"
	"os"
	"strconv"
)

type FlagsResult struct {
	Host string
	Port uint16

	LedCount int
	FftSize  int

	SampleRate int
	Channels   int

	Decay         float64
	DbfsThreshold float64

	Backend string
	Device  string

	AudibleLow  float64
	AudibleHigh float64

	Mirror bool
	Effect string

	Color color.RGBA

	Verbose bool
}

func GetFlags() FlagsResult {
	var host = flag.String("host", "", "host of the luxsrv")
	var port = flag.Uint("port", DefaultPort, "port of the luxsrv")

	var ledCount = flag.Int("leds", 0, "number of LEDs to be driven (max 255)")
	var fftSize = flag.Int("fft", 1024, "FFT size")

	var sampleRate = flag.Int("sampleRate", 0, "sample rate")
	var channels = flag.Int("channels", 2, "number of channels")

	var decay = flag.Float64("decay", 0.50, "decay factor [0,1] controls the smoothness of the visualization")
	var dbfsThreshold = flag.Float64("dbfsThreshold", -GetSQNR(16), "dBFS threshold")

	var backend = flag.String("backend", "auto", "audio backend (auto, wasapi, alsa, pulse, jack)")
	var device = flag.String("device", "loopback", "device to use (loopback, capture)")

	var audibleLow = flag.Float64("audibleLow", 20, "lower audible frequency")
	var audibleHigh = flag.Float64("audibleHigh", 20000, "upper audible frequency")

	var mirror = flag.Bool("mirror", false, "mirror mode with lower frequencies at the middle")
	var effect = flag.String("effect", "solid", "color effect (solid, rainbow, luxception)")

	var color = flag.String("color", "ff00ff", "hex color")

	var verbose = flag.Bool("verbose", false, "print verbose messages")

	flag.Parse()

	if *host == "" || *ledCount == 0 || *ledCount > 255 || *sampleRate == 0 {
		flag.Usage()
		os.Exit(2)
	}

	rgb, err := parseColor(*color)
	if err != nil {
		flag.Usage()
		os.Exit(2)
	}

	return FlagsResult{
		Host: *host,
		Port: uint16(*port),

		LedCount: *ledCount,
		FftSize:  *fftSize,

		SampleRate: *sampleRate,
		Channels:   *channels,

		Decay:         *decay,
		DbfsThreshold: *dbfsThreshold,

		Backend: *backend,
		Device:  *device,

		AudibleLow:  *audibleLow,
		AudibleHigh: *audibleHigh,

		Mirror: *mirror,
		Effect: *effect,

		Color: rgb,

		Verbose: *verbose,
	}
}

func parseColor(s string) (rgb color.RGBA, err error) {
	c, err := strconv.ParseUint(s, 16, 24)

	rgb.R = uint8((c & 0xFF0000) >> 16)
	rgb.G = uint8((c & 0x00FF00) >> 8)
	rgb.B = uint8((c & 0x0000FF))

	return rgb, err
}
