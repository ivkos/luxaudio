package analyzers

type Analyzer interface {
	Analyze([]float64) []byte
}
