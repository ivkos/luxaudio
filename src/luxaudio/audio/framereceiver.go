package audio

import (
	"bytes"
	"encoding/binary"
	"luxaudio/analyzers"
	"luxaudio/utils"
)

type FrameReceiver struct {
	sampleSizeInBytes int
	channels          int
	queue             *analyzers.Queue
}

func NewFrameReceiver(sampleSizeInBytes int, channels int, queue *analyzers.Queue) *FrameReceiver {
	return &FrameReceiver{
		sampleSizeInBytes: sampleSizeInBytes,
		channels:          channels,
		queue:             queue,
	}
}

func (fr *FrameReceiver) OnReceive(data []byte, frameCount uint32) {
	convertedData := make([]float32, len(data)/fr.sampleSizeInBytes)
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &convertedData)
	utils.CheckErr(err)

	// downsample to mono
	monoFloats := fr.downsampleToMono(convertedData)

	fr.queue.Enqueue(monoFloats, false)
}

func (fr *FrameReceiver) downsampleToMono(convertedData []float32) []float64 {
	monoFloats := make([]float64, len(convertedData)/fr.channels)

	for i := range monoFloats {
		for j := 0; j < fr.channels; j++ {
			monoFloats[i] += float64(convertedData[i*fr.channels+j])
		}
		monoFloats[i] = monoFloats[i] / float64(fr.channels)
	}

	return monoFloats
}
