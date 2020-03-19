package solver

//
//import (
//	"libRPSO/vector"
//	"math"
//	"math/rand"
//)
//
//type Solution struct {
//	position  []float64
//	evalValue float64
//	velocity  []float64
//}
//
//func NewSolution(dim int) *Solution {
//	return &Solution{
//		position:  make([]float64, dim),
//		velocity:  make([]float64, dim),
//		evalValue: math.MaxFloat64,
//	}
//}
//
//func (s *Solution) update(x []float64, e float64) {
//	s.position = x
//	s.evalValue = e
//}
//
//func (s *Solution) copy(src *Solution) {
//	copy(s.position, src.position)
//	s.evalValue = src.evalValue
//}
//
//func (s *Solution) GetEval() float64 {
//	return s.evalValue
//}
//
//func (s *Solution) GetPosition() []float64 {
//	return s.position
//}
//
//func (s *Solution) learn(px, py []float64, c float64) {
//	tmp := vector.Add(px, py, 1, -1)
//	r := rand.Float64()
//	tmp = vector.Scale(tmp, c*r)
//	s.velocity = vector.Add(s.velocity, tmp, 1, 1)
//}
//
//func DefaultInitSolution(bound *Bound, dim int) *Solution {
//	s := NewSolution(dim)
//	for i := 0; i < dim; i++ {
//		s.position[i] = rangeRand(bound.XUpper, bound.XLower)
//		s.velocity[i] = rangeRand(bound.VUpper, bound.VLower)
//	}
//	return s
//}
//
//func rangeRand(upper, lower float64) float64 {
//	return (upper-lower)*rand.Float64() + lower
//}
