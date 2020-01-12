package utils

import (
	"gonum.org/v1/gonum/stat"
)

func Chunk(data []float64, desiredChunks int) [][]float64 {
	var divided [][]float64

	chunkSize := (len(data) + desiredChunks - 1) / desiredChunks

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize

		if end > len(data) {
			end = len(data)
		}

		divided = append(divided, data[i:end])
	}

	return divided
}

func ChunkedMean(data []float64, desiredChunks int) []float64 {
	chunks := Chunk(data, desiredChunks)

	result := make([]float64, len(chunks))
	for i := range result {
		result[i] = stat.Mean(chunks[i], nil)
	}

	return result
}

func CenterArray(arr []float64, total int) []float64 {
	arrLen := len(arr)
	if total <= arrLen {
		return arr
	}

	offset := (total - arrLen) / 2

	result := make([]float64, 0)
	result = append(result, make([]float64, offset)...)
	result = append(result, arr...)
	result = append(result, make([]float64, total-len(result))...)

	return result
}
