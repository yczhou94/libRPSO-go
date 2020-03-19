package solver

import (
	"math"
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
}

func NewSolverConf(maxStep int) *Conf {
	return &Conf{
		MaxStep:     maxStep,
		NTerm:       maxStep,
		PrintEvery:  1,
		Convergence: 1e-5,
	}
}

func NewSolver(psoParam *PSOParam, conf *Conf) (*Solver, error) {
	var err error
	solver := &Solver{
		conf: conf,
	}
	solver.swarm, err = newSwarm(psoParam)

	return solver, err
}

func (s *Solver) Run() error {
	var tCount int
	for i := 0; i < s.conf.MaxStep; i++ {
		old := s.swarm.gBestMem.evalValue
		err := s.swarm.step()
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

	return nil
}
