# luxaudio
❤️💚💙

**luxaudio** is a Go application that captures audio, does spectral analysis 
on it, and sends a spectrum visualization to a **[luxsrv](https://github.com/ivkos/luxsrv)**-enabled RGB LED strip 
over the network.

**luxaudio** is part of **[Lux](https://github.com/ivkos/lux)**.


## Requirements
* Go

## Quick Start

### Install dependencies
`GOPATH=$PWD go get -v -d luxaudio`

### Build
`GOPATH=$PWD go build -v luxaudio`

### Start
```
./luxaudio \
    --channels 1 \
    --device capture \
    --host 10.10.10.108 \
    --leds 120 \
    --sampleRate 44100 \
    --fft 1024 \
    --decay 0.5 \
    --dbfsThreshold -64 \
    --audibleLow 30 \
    --audibleHigh 17000 \
    --mirror true
```

### Usage
```
Usage of ./luxaudio:
  -audibleHigh float
        upper audible frequency (default 20000)
  -audibleLow float
        lower audible frequency (default 20)
  -backend string
        audio backend (auto, wasapi, alsa, pulse, jack) (default "auto")
  -channels int
        number of channels (default 2)
  -dbfsThreshold float
        dBFS threshold (default -96.32959861247399)
  -decay float
        decay factor [0,1] controls the smoothness of the visualization (default 0.5)
  -device string
        device to use (loopback, capture) (default "loopback")
  -fft int
        FFT size (default 1024)
  -host string
        host of the luxsrv
  -leds int
        number of LEDs to be driven (max 255)
  -mirror
        mirror mode with lower frequencies at the middle
  -port uint
        port of the luxsrv (default 42170)
  -sampleRate int
        sample rate
```
