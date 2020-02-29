package analyzers

type Analyzer interface {
	Analyze([]float64) []float64
}
