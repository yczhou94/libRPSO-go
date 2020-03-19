package solver

import (
	"log"
	"math"
	"math/rand"
	"time"
)

type TargetFunc func(x []float64, args ...interface{}) ([]float64, float64, error)

type Solver struct {
	swarm *swarm
	conf  *Conf
}

type Conf struct {
	MaxStep     int
	NTerm       int
	PrintEvery  int
	Convergence float64
	Seed        int64
}

func NewSolverConf(maxStep int) *Conf {
	return &Conf{
		MaxStep:     maxStep,
		NTerm:       maxStep,
		PrintEvery:  1,
		Convergence: 1e-5,
		Seed:        time.Now().Unix(),
	}
}

func NewSolver(psoParam *PSOParam, conf *Conf) (*Solver, error) {
	var err error
	solver := &Solver{
		conf: conf,
	}
	psoParam.maxStep = conf.MaxStep
	solver.swarm, err = newSwarm(psoParam)
	rand.Seed(conf.Seed)

	return solver, err
}

func (s *Solver) Run() error {
	var (
		tCount int
		step   int
	)
	for step = 0; step < s.conf.MaxStep; step++ {
		old := s.swarm.gBestMem.evalValue
		err := s.swarm.step(step)
		if err != nil {
			return err
		}
		if math.Abs(old-s.swarm.gBestMem.evalValue) < s.conf.Convergence {
			tCount++
		} else {
			tCount = 0
		}
		if tCount == s.conf.NTerm {
			break
		}
	}
	s.finish(step)
	return nil
}

func (s *Solver) finish(step int) {
	log.Printf("RPSO stop at step %d\n", step)
	solution := s.swarm.getBestSolution()
	log.Printf("final evaluate value of target function: %f\n", solution.evalValue)
	log.Printf("final coordinates of optimal solution: %+v\n", solution.position)
}

func (s *Solver) GetSolution() *Solution {
	return s.swarm.getBestSolution()
}
