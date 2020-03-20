package solver

import (
	"libRPSO/vector"
	"math"
	"math/rand"
)

type Bound struct {
	XUpper float64
	XLower float64
	VUpper float64
	VLower float64
}

type Solution struct {
	position  []float64
	evalValue float64
	velocity  []float64
}

type InitSolutionFunc func(bound *Bound, dim int) *Solution

func NewSolution(param *PSOParam) (*Solution, error) {
	var err error
	s := param.initFunc(param.bound, param.dim)
	err = s.setUp(param)
	return s, err
}

func (s *Solution) setUp(param *PSOParam) error {
	x, e, err := param.targetFunc(s.position, param.args)
	if err != nil {
		return err
	}
	s.update(x, e)
	return nil
}

func (s *Solution) update(x []float64, e float64) {
	s.position = x
	s.evalValue = e
}

func (s *Solution) Copy(src *Solution) {
	copy(s.position, src.position)
	s.evalValue = src.evalValue
}

func (s *Solution) GetEval() float64 {
	return s.evalValue
}

func (s *Solution) GetPosition() []float64 {
	return s.position
}

func (s *Solution) learn(px, py []float64, c float64) {
	tmp := vector.Add(px, py, 1, -1)
	r := rand.Float64()
	tmp = vector.Scale(tmp, c*r)
	s.velocity = vector.Add(s.velocity, tmp, 1, 1)
}

func DefaultInitSolution(bound *Bound, dim int) *Solution {
	s := &Solution{
		position:  make([]float64, dim),
		velocity:  make([]float64, dim),
		evalValue: math.MaxFloat64,
	}
	for i := 0; i < dim; i++ {
		s.position[i] = rangeRand(bound.XUpper, bound.XLower)
		s.velocity[i] = rangeRand(bound.VUpper, bound.VLower)
	}
	return s
}

func rangeRand(upper, lower float64) float64 {
	return (upper-lower)*rand.Float64() + lower
}

func (s *Solution) step(idx int, pBest []*Solution, gBest *Solution, param *PSOParam) error {
	s.position = vector.Add(s.position, s.velocity, 1, 1)
	s.updateVelocity(idx, pBest, gBest, param)

	solution := &Solution{}
	x := s.position
	pm := rand.Float64()
	if pm < param.pm {
		solution = param.initFunc(param.bound, param.dim)
		x = solution.position
	}

	x, e, err := param.targetFunc(x, param.args)
	if err != nil {
		return err
	}
	s.update(x, e)

	return nil
}

func (s *Solution) updateVelocity(idx int, pBest []*Solution, gBest *Solution, param *PSOParam) {
	// v = w*v
	s.velocity = vector.Scale(s.velocity, param.w)
	// v += c1*r1*(pBest[i] - pos)
	s.learn(pBest[idx].position, s.position, param.c1)
	// v += c2*r2*(gBest - pos)
	s.learn(gBest.position, s.position, param.c2)

	// v += c3*r3*(pBest[r] - pos)
	pr := rand.Float64()
	if pr < param.pr {
		rIdx := rand.Intn(len(pBest))
		s.learn(pBest[rIdx].position, s.position, param.c3)
	}
}
