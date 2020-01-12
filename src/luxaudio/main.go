package main

import (
	"fmt"
	"github.com/gen2brain/malgo"
	"log"
	"luxaudio/analyzers"
	"luxaudio/audio"
	"luxaudio/led"
	"luxaudio/utils"
	"runtime"
	"time"
)

func main() {
	f := utils.GetFlags()

	malgoBackend := getBackend(f.Backend)
	malgoDevice := getDevice(f.Device)

	context, captureConfig := initMalgo(uint32(f.Channels), uint32(f.SampleRate), malgoBackend, malgoDevice)
	defer func() {
		_ = context.Uninit()
		context.Free()
	}()

	// Create UDP socket
	conn := utils.GetUDPConn(f.Host, f.Port)
	defer func() { _ = conn.Close() }()

	payloadSender := func(ledData []byte) {
		_, err := conn.Write(led.MakeRawModeLuxPayload(uint8(f.LedCount), ledData))
		utils.CheckErr(err)
	}

	analyzer := analyzers.NewSmartAnalyzer(
		f.FftSize,
		f.LedCount,
		float64(f.SampleRate),
		f.Decay,
		f.DbfsThreshold,
		f.AudibleLow,
		f.AudibleHigh,
	)

	queue := analyzers.NewQueue(f.FftSize, &analyzer, &payloadSender)
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

func initMalgo(channels uint32, sampleRate uint32, backend malgo.Backend, device malgo.DeviceType) (*malgo.AllocatedContext, malgo.DeviceConfig) {
	ctxConfig := malgo.ContextConfig{}
	ctxConfig.ThreadPriority = malgo.ThreadPriorityRealtime

	context, err := malgo.InitContext([]malgo.Backend{backend}, ctxConfig, func(message string) {
		log.Printf("LOG <%v>\n", message)
	})
	utils.CheckErr(err)

	captureConfig := malgo.DefaultDeviceConfig()
	captureConfig.PerformanceProfile = malgo.LowLatency
	captureConfig.DeviceType = device
	captureConfig.Capture.Format = malgo.FormatF32
	captureConfig.SampleRate = sampleRate
	captureConfig.Capture.Channels = channels

	return context, captureConfig
}

func getBackend(backend string) malgo.Backend {
	switch backend {
	case "auto":
		switch os := runtime.GOOS; os {
		case "linux":
			return malgo.BackendAlsa

		case "windows":
			return malgo.BackendWasapi

		default:
			log.Fatalf("Unsupported operating system: %s", os)
		}

	case "alsa":
		return malgo.BackendAlsa

	case "pulse":
		return malgo.BackendPulseaudio

	case "jack":
		return malgo.BackendJack

	case "wasapi":
		return malgo.BackendWasapi

	default:
		log.Fatalf("Unsupported backend: %s", backend)
	}

	return malgo.BackendNull
}

func getDevice(device string) malgo.DeviceType {
	switch device {
	case "loopback":
		return malgo.Loopback

	case "capture":
		return malgo.Capture

	default:
		log.Fatalf("Unsupported device: %s", device)
	}

	return 0
}
