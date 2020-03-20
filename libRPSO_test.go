package libRPSO

import (
	"fmt"
	"libRPSO/solver"
	"libRPSO/vector"
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

func distance(x, y []float64) float64 {
	vector.CheckEqualLen(x, y)
	var dist float64
	for i := range x {
		dist += (x[i] - y[i]) * (x[i] - y[i])
	}
	return dist
}

func TestSolver_Run(t *testing.T) {
	popSize := 100
	dim := 100
	maxStep := 20000
	psoParam := solver.NewPSOParam(popSize, dim, Rosenbrock)
	psoParam.SetNProc(8)

	psoParam.SetBound(&solver.Bound{
		XUpper: 30,
		XLower: -30,
		VUpper: 1,
		VLower: -1,
	})

	times := 10
	eval := make([]float64, 7)

	for i := 0; i < times; i++ {
		log.Println(i)
		conf := solver.NewSolverConf(maxStep)
		conf.Seed = time.Now().Unix()
		conf.NTerm = 2000

		p := make([]*solver.Solution, popSize)
		for i := 0; i < popSize; i++ {
			p[i] = solver.DefaultInitSolution(psoParam.GetBound(), dim)
		}

		pMem := make([]*solver.Solution, popSize)
		for i := range pMem {
			pMem[i], _ = solver.NewSolution(psoParam)
		}
		//copySolutions(pMem, p)

		//log.Println("native")
		//psoParam.SetSimAnnealFlag(false)
		//psoParam.SetPm(0)
		//psoParam.SetPr(0)
		//conf.NConfusion = 0
		//eval[0] += runSolver(psoParam, conf, p)
		//
		//log.Println("random")
		//copySolutions(p, pMem)
		//psoParam.SetPr(1)
		//psoParam.SetW(0)
		//psoParam.SetC1(0)
		//psoParam.SetC2(0)
		//eval[1] += runSolver(psoParam, conf, p)
		//
		//log.Println("sim anneal")
		//copySolutions(p, pMem)
		//psoParam.SetSimAnnealFlag(true)
		//psoParam.SetT(10000.0)
		//eval[2] += runSolver(psoParam, conf, p)

		//log.Println("native")
		//psoParam.SetSimAnnealFlag(false)
		//psoParam.SetPr(0)
		//psoParam.SetPm(0)
		//conf.NConfusion = 0
		//eval[0] += runSolver(psoParam, conf, p)

		//copySolutions(p, pMem)
		//log.Println("sim anneal")
		//psoParam.SetSimAnnealFlag(true)
		//eval[1] += runSolver(psoParam, conf, p)
		//
		//copySolutions(p, pMem)
		//log.Println("random learning")
		//psoParam.SetSimAnnealFlag(false)
		//psoParam.SetPr(0.1)
		//eval[2] += runSolver(psoParam, conf, p)
		//
		//copySolutions(p, pMem)
		//log.Println("mutate")
		//psoParam.SetPr(0)
		//psoParam.SetPm(0.01)
		//eval[3] += runSolver(psoParam, conf, p)
		//
		//copySolutions(p, pMem)
		//log.Println("confusion")
		//conf.NConfusion = int(math.Sqrt(float64(maxStep)))
		//psoParam.SetPm(0)
		//eval[4] += runSolver(psoParam, conf, p)

		log.Println("revised")
		copySolutions(p, pMem)
		//p = genSolWithMinDist(popSize, dim, psoParam)
		psoParam.SetPr(0.1)
		psoParam.SetPm(0.01)
		psoParam.SetSimAnnealFlag(true)
		//psoParam.SetT(1000.0)
		eval[5] += runSolver(psoParam, conf, p)

		//p = genSolWithMinDist(popSize, dim, psoParam)
		//eval[6] += runSolver(psoParam, conf, p)
	}

	fmt.Printf("MSE: %+v\n", vector.Scale(eval, 1.0/float64(times)))
}

func genSolWithMinDist(popSize, dim int, param *solver.PSOParam) []*solver.Solution {
	var solutions []*solver.Solution
	for len(solutions) < popSize {
		sol := solver.DefaultInitSolution(param.GetBound(), dim)
		if len(solutions) > 0 {
			dist := getMinDistanceInSwarm(sol, solutions)
			if dist > 2.0 {
				solutions = append(solutions, sol)
			}
		} else {
			solutions = append(solutions, sol)
		}
	}

	return solutions
}

func copySolutions(dst, src []*solver.Solution) {
	for i := range src {
		dst[i].Copy(src[i])
	}
}

func getMinDistanceInSwarm(sol *solver.Solution, solutions []*solver.Solution) float64 {
	dist := math.MaxFloat64
	for _, s := range solutions {
		tmp := distance(sol.GetPosition(), s.GetPosition())
		if tmp < dist {
			dist = tmp
		}
	}
	return dist
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
