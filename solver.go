package libRPSO

type TargetFunc func(x []float64, args ...interface{}) ([]float64, float64, error)
