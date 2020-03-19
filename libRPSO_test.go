package libRPSO

import (
	"fmt"
	"libRPSO/solver"
	"log"
	"math"
	"testing"
	"time"
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

func Rosenbrock(x []float64, args ...interface{}) ([]float64, float64, error) {
	if args == nil {

	}
	var sum float64
	for i := 0; i < len(x)-1; i++ {
		sum += 100*math.Pow(x[i+1]-x[i], 2) + math.Pow(1-x[i], 2)
	}
	return x, sum, nil
}

func TestSolver_Run(t *testing.T) {
	popSize := 50
	dim := 50
	maxStep := 10000
	psoParam := solver.NewPSOParam(popSize, dim, Rosenbrock)
	psoParam.SetNProc(8)

	times := 10
	eval := make([]float64, 6)

	for i := 0; i < times; i++ {
		conf := solver.NewSolverConf(maxStep)
		conf.Seed = time.Now().Unix()

		p := make([]*solver.Solution, popSize)
		for i := 0; i < popSize; i++ {
			p[i] = solver.DefaultInitSolution(psoParam.GetBound(), dim)
		}

		pMem := make([]*solver.Solution, popSize)
		copy(pMem, p)

		log.Println("native")
		psoParam.SetSimAnnealFlag(false)
		psoParam.SetPr(0)
		psoParam.SetPm(0)
		conf.NConfusion = 0
		eval[0] += runSolver(psoParam, conf, p)

		copy(p, pMem)
		log.Println("sim anneal")
		psoParam.SetSimAnnealFlag(true)
		eval[1] += runSolver(psoParam, conf, p)

		copy(p, pMem)
		log.Println("random learning")
		psoParam.SetSimAnnealFlag(false)
		psoParam.SetPr(0.1)
		eval[2] += runSolver(psoParam, conf, p)

		copy(p, pMem)
		log.Println("mutate")
		psoParam.SetPr(0)
		psoParam.SetPm(0.01)
		eval[3] += runSolver(psoParam, conf, p)

		copy(p, pMem)
		log.Println("confusion")
		conf.NConfusion = int(math.Sqrt(float64(maxStep)))
		psoParam.SetPm(0)
		eval[4] += runSolver(psoParam, conf, p)

		copy(p, pMem)
		log.Println("revised")
		psoParam.SetPr(0.1)
		psoParam.SetPm(0.01)
		psoParam.SetSimAnnealFlag(true)
		eval[5] += runSolver(psoParam, conf, p)
	}

	fmt.Printf("%+v\n", eval)
}

func runSolver(param *solver.PSOParam, conf *solver.Conf, solutions []*solver.Solution) float64 {
	s, err := solver.NewSolver(param, conf, solutions)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}

	return s.GetSolution().GetEval()
}
