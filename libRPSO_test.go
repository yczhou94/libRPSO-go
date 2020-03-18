package libRPSO

import (
	"math/rand"
	"testing"
)

func initParticle(bound *Bound, dim int) (*Solution, *Velocity) {
	s := NewSolution(dim)
	v := NewVelocity(dim)
	for i := 0; i < dim; i++ {
		s.position[i] = rangeRand(bound.XUpper, bound.XLower)
		v.v[i] = rangeRand(bound.VUpper, bound.VLower)
	}
	return s, v
}

func rangeRand(upper, lower float64) float64 {
	return (upper-lower)*rand.Float64() + lower
}

func SumSquare(x []float64, args ...interface{}) ([]float64, float64, error) {
	var sum float64
	for i := 0; i < len(x); i++ {
		sum += x[i] * x[i]
	}
	return x, sum, nil
}

func TestSolver_Run(t *testing.T) {
	psoParam := NewPSOParam()
	psoParam.Dim = 10
	psoParam.PopSize = 10
	psoParam.InitFunc = initParticle
	psoParam.Target = SumSquare
	psoParam.NProc = 4

	conf := &SolverConf{
		MaxStep:     100,
		NTerm:       100,
		PrintEvery:  1,
		Convergence: 0.01,
	}
	solver, err := NewSolver(psoParam, conf)
	if err != nil {
		t.Fatal(err)
	}
	err = solver.Run()
	if err != nil {
		t.Fatal(err)
	}
}
