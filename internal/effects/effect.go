package effects

type Effect interface {
	Apply([]float64) []byte
}
