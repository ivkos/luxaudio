package utils

import (
	"log"
	"math"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetSQNR(bits int) float64 {
	return 20 * math.Log10(math.Pow(2, float64(bits)))
}
