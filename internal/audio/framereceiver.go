package audio

import (
	"bytes"
	"encoding/binary"
	"github.com/ivkos/luxaudio/internal/analyzers"
	"github.com/ivkos/luxaudio/internal/utils"
)

type FrameReceiver struct {
	sampleSizeInBytes int
	channels          int
	queue             *analyzers.Queue
	pinger            *utils.Pinger
}

type SampleFormat float32

func NewFrameReceiver(sampleSizeInBytes int, channels int, queue *analyzers.Queue, pinger *utils.Pinger) *FrameReceiver {
	return &FrameReceiver{
		sampleSizeInBytes: sampleSizeInBytes,
		channels:          channels,
		queue:             queue,
		pinger:            pinger,
	}
}

func (fr *FrameReceiver) OnReceive(data []byte, frameCount uint32) {
	if !fr.pinger.IsReachable {
		return
	}

	convertedData := make([]SampleFormat, len(data)/fr.sampleSizeInBytes)
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &convertedData)
	utils.CheckErr(err)

	// downsample to mono
	monoFloats := fr.downsampleToMono(convertedData)

	fr.queue.Enqueue(monoFloats, false)
}

func (fr *FrameReceiver) downsampleToMono(convertedData []SampleFormat) []float64 {
	monoFloats := make([]float64, len(convertedData)/fr.channels)

	for i := range monoFloats {
		for j := 0; j < fr.channels; j++ {
			monoFloats[i] += float64(convertedData[i*fr.channels+j])
		}
		monoFloats[i] = monoFloats[i] / float64(fr.channels)
	}

	return monoFloats
}
