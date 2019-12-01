package analyzers

import (
	"log"
)

type PayloadSender = func(ledData []byte)

type Queue struct {
	fftSize     int
	analyzer    *Analyzer
	sampleQueue []float64

	sender *PayloadSender
}

func NewQueue(fftSize int, analyzer *Analyzer, sender *PayloadSender) *Queue {
	return &Queue{
		fftSize:     fftSize,
		analyzer:    analyzer,
		sampleQueue: make([]float64, 0),
		sender:      sender,
	}
}

func (q *Queue) Size() int {
	return len(q.sampleQueue)
}

func (q *Queue) Enqueue(monoFloats []float64, recursiveCall bool) {
	q.sampleQueue = append(q.sampleQueue, monoFloats...)

	if len(q.sampleQueue) < q.fftSize {
		return
	}

	if recursiveCall {
		log.Printf("Leftover samples")
	}

	// get our chunk
	sampleChunk := q.sampleQueue[:q.fftSize]

	// analyze
	ledData := (*(q.analyzer)).Analyze(sampleChunk)

	// remove analyzed chunk
	q.sampleQueue = q.sampleQueue[q.fftSize:]

	// send the payload
	(*(q.sender))(ledData)

	q.Enqueue([]float64{}, true)
}
