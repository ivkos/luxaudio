package utils

import "math"

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
		result[i] = MeanOfSlice(chunks[i])
	}

	return result
}

func MeanOfSlice(data []float64) float64 {
	dataLen := len(data)
	if dataLen == 0 {
		return 0
	}

	var total float64 = 0

	for _, value := range data {
		total += value
	}

	return total / float64(dataLen)
}

func RmsOfSlice(data []float64) float64 {
	var total float64 = 0

	for _, value := range data {
		total += value * value
	}

	return math.Sqrt(total / float64(len(data)))
}

func MaxOfSlice(data []float64) float64 {
	max := math.Inf(-1)

	for _, x := range data {
		if x > max {
			max = x
		}
	}

	return max
}

func MaxAbsOfSlice(data []float64) float64 {
	max := math.Inf(-1)

	for _, x := range data {
		abs := math.Abs(x)
		if abs > max {
			max = abs
		}
	}

	return max
}

func SumOfSlice(data []float64) float64 {
	sum := float64(0)

	for _, x := range data {
		sum += x
	}

	return sum
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
