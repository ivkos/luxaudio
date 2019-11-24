package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/gen2brain/malgo"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
	"luxaudio/led"
	"luxaudio/utils"
	"math/cmplx"
	"os"
)

const DefaultPort = 42170

func main() {
	luxsrvHost, luxsrvPort, ledCount, fftSize := getFlags()

	context, captureConfig := initMalgo()
	defer func() {
		_ = context.Uninit()
		context.Free()
	}()

	// Create UDP socket
	conn := utils.GetUDPConn(luxsrvHost, luxsrvPort)
	defer func() { _ = conn.Close() }()

	ledData := make([]byte, ledCount*3)
	magnitudes := make([]float64, fftSize/2)
	sampleChan := make(chan float64, fftSize)
	dataFromChannel := make([]float64, fftSize)
	sizeInBytes := malgo.SampleSizeInBytes(captureConfig.Capture.Format)

	onReceiveFrames := func(_, data []byte, frameCount uint32) {
		dataLen := len(data)

		actualData := make([]int16, dataLen/sizeInBytes)
		err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &actualData)
		utils.CheckErr(err)

		monoFloats := make([]float64, dataLen/sizeInBytes/2)

		for i := range monoFloats {
			average := (float64(actualData[i]) + float64(actualData[i+1])) / 2
			monoFloats[i] = average

			if len(sampleChan) >= fftSize {
				<-sampleChan
			}

			sampleChan <- monoFloats[i]
		}

		if len(sampleChan) < fftSize {
			return
		}

		for i := 0; i < fftSize; i++ {
			dataFromChannel[i] = <-sampleChan
		}
		window.Apply(dataFromChannel, window.Hann)

		ffs := fft.FFTReal(dataFromChannel)
		for i := range magnitudes {
			magnitudes[i] = cmplx.Abs(ffs[i])
		}

		for i := 0; i < ledCount; i++ {
			v := float64(0)
			for j := 0; j < 2; j++ {
				v = v + magnitudes[2*i+j]
			}
			v = v / float64(fftSize)
			if v > 255 {
				ledData[3*i+0] = 255
				ledData[3*i+1] = 0
				ledData[3*i+2] = 0
			} else {
				ledData[3*i+0] = byte(v)
				ledData[3*i+1] = 0
				ledData[3*i+2] = 0
			}
		}

		_, err = conn.Write(led.MakeRawModeLuxPayload(uint8(ledCount), ledData))
		utils.CheckErr(err)
	}

	fmt.Println("Listening...")
	device, err := malgo.InitDevice(context.Context, captureConfig, malgo.DeviceCallbacks{
		Data: onReceiveFrames,
	})
	utils.CheckErr(err)

	defer device.Uninit()

	err = device.Start()
	utils.CheckErr(err)

	fmt.Println("Press Enter to stop listening...")
	fmt.Scanln()
}

func getFlags() (string, uint16, int, int) {
	var host = flag.String("host", "", "host of the luxsrv")
	var port = flag.Uint("port", DefaultPort, "port of the luxsrv")
	var ledCount = flag.Int("leds", 0, "number of LEDs to be driven (max 255)")
	var fftSize = flag.Int("fft", 2048, "FFT size")

	flag.Parse()

	if *host == "" || *ledCount == 0 || *ledCount > 255 {
		flag.Usage()
		os.Exit(2)
	}

	return *host, uint16(*port), *ledCount, *fftSize
}

func initMalgo() (*malgo.AllocatedContext, malgo.DeviceConfig) {
	ctxConfig := malgo.ContextConfig{}
	ctxConfig.ThreadPriority = malgo.ThreadPriorityRealtime

	context, err := malgo.InitContext([]malgo.Backend{malgo.BackendWasapi}, ctxConfig, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})

	utils.CheckErr(err)

	captureConfig := malgo.DefaultDeviceConfig()
	captureConfig.PerformanceProfile = malgo.LowLatency
	captureConfig.DeviceType = malgo.Loopback
	captureConfig.Capture.Format = malgo.FormatS16
	// captureConfig.SampleRate = 96000
	captureConfig.Capture.Channels = 2

	return context, captureConfig
}
