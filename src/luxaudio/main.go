package main

import (
	"fmt"
	"github.com/gen2brain/malgo"
	"log"
	"luxaudio/analyzers"
	"luxaudio/audio"
	"luxaudio/led"
	"luxaudio/utils"
	"time"
)

func main() {
	luxsrvHost, luxsrvPort, ledCount, fftSize, sampleRate, channels, decayFactor := utils.GetFlags()

	context, captureConfig := initMalgo(uint32(channels), uint32(sampleRate))
	defer func() {
		_ = context.Uninit()
		context.Free()
	}()

	// Create UDP socket
	conn := utils.GetUDPConn(luxsrvHost, luxsrvPort)
	defer func() { _ = conn.Close() }()

	payloadSender := func(ledData []byte) {
		_, err := conn.Write(led.MakeRawModeLuxPayload(uint8(ledCount), ledData))
		utils.CheckErr(err)
	}

	analyzer := analyzers.NewSmartAnalyzer(fftSize, ledCount, float64(sampleRate), decayFactor)
	queue := analyzers.NewQueue(fftSize, &analyzer, &payloadSender)
	frameReceiver := audio.NewFrameReceiver(
		malgo.SampleSizeInBytes(captureConfig.Capture.Format),
		int(captureConfig.Capture.Channels),
		queue,
	)

	go func() {
		for {
			t := time.NewTimer(1 * time.Second)
			<-t.C
			log.Printf("len(queue) = %d\n", queue.Size())
		}
	}()

	log.Println("Listening...")
	device, err := malgo.InitDevice(context.Context, captureConfig, malgo.DeviceCallbacks{
		Data: func(_, data []byte, count uint32) { frameReceiver.OnReceive(data, count) },
	})
	utils.CheckErr(err)

	defer device.Uninit()

	err = device.Start()
	utils.CheckErr(err)

	log.Println("Press Enter to stop listening...")
	fmt.Scanln()
}

func initMalgo(channels uint32, sampleRate uint32) (*malgo.AllocatedContext, malgo.DeviceConfig) {
	ctxConfig := malgo.ContextConfig{}
	ctxConfig.ThreadPriority = malgo.ThreadPriorityRealtime

	context, err := malgo.InitContext([]malgo.Backend{malgo.BackendWasapi}, ctxConfig, func(message string) {
		log.Printf("LOG <%v>\n", message)
	})
	utils.CheckErr(err)

	captureConfig := malgo.DefaultDeviceConfig()
	captureConfig.PerformanceProfile = malgo.LowLatency
	captureConfig.DeviceType = malgo.Loopback
	captureConfig.Capture.Format = malgo.FormatF32
	captureConfig.SampleRate = sampleRate
	captureConfig.Capture.Channels = channels

	return context, captureConfig
}
