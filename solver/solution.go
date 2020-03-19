package solver

import "math"

type Solution struct {
	position  []float64
	evalValue float64
}

func NewSolution(dim int) *Solution {
	return &Solution{
		position:  make([]float64, dim),
		evalValue: math.MaxFloat64,
	}
}

func (s *Solution) update(x []float64, e float64) {
	s.position = x
	s.evalValue = e
}

func (s *Solution) copy(src *Solution) {
	copy(s.position, src.position)
	s.evalValue = src.evalValue
}
