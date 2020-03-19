package libRPSO

import (
	"libRPSO/solver"
	"testing"
)

func SumSquare(x []float64, args ...interface{}) ([]float64, float64, error) {
	if args == nil {

	}
	var sum float64
	for i := 0; i < len(x); i++ {
		sum += x[i] * x[i]
	}
	return x, sum, nil
}

func TestSolver_Run(t *testing.T) {
	psoParam := solver.NewPSOParam(10, 10, SumSquare)
	psoParam.SetNProc(4)

	conf := solver.NewSolverConf(100)

	s, err := solver.NewSolver(psoParam, conf)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		t.Fatal(err)
	}
}
