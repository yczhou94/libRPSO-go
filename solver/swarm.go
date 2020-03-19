package solver

import (
	"math"
	"math/rand"
	"sync"
)

type swarm struct {
	popSize   int
	solutions []*Solution
	pBest     []*Solution
	gBest     *Solution
	pBestMem  []*Solution
	gBestMem  *Solution
	param     *PSOParam
	wg        sync.WaitGroup
}

func newSwarm(param *PSOParam, solutions []*Solution) (*swarm, error) {
	var err error
	s := &swarm{
		popSize:   param.popSize,
		solutions: make([]*Solution, param.popSize),
		pBest:     make([]*Solution, param.popSize),
		pBestMem:  make([]*Solution, param.popSize),
		param:     param,
	}

	s.gBest, err = NewSolution(param)
	if err != nil {
		return nil, err
	}
	s.gBestMem, err = NewSolution(param)
	if err != nil {
		return nil, err
	}

	if solutions == nil {
		solutions = make([]*Solution, param.popSize)
		for i := 0; i < param.popSize; i++ {
			solutions[i], err = NewSolution(param)
			if err != nil {
				return nil, err
			}
		}
	} else {
		for i := 0; i < param.popSize; i++ {
			err = solutions[i].setUp(param)
			if err != nil {
				return nil, err
			}
		}
	}

	s.solutions = solutions

	for i := 0; i < param.popSize; i++ {
		s.pBest[i], err = NewSolution(param)
		if err != nil {
			return nil, err
		}
		s.pBestMem[i], err = NewSolution(param)
		if err != nil {
			return nil, err
		}
		s.pBest[i].copy(s.solutions[i])
	}

	copy(s.pBestMem, s.pBest)
	s.updateBest()

	return s, nil
}

func (s *swarm) step() error {
	var err error
	s.wg.Add(s.param.popSize)
	queue := make(chan struct{}, s.param.nProc)
	for idx, _ := range s.solutions {
		queue <- struct{}{}
		worker := func(idx int) {
			defer func() {
				<-queue
				s.wg.Done()
			}()
			solution := s.solutions[idx]
			e := solution.step(idx, s.pBest, s.gBest, s.param)
			if e != nil {
				err = e
				return
			}
			if solution.evalValue < s.pBest[idx].evalValue {
				s.pBest[idx].copy(solution)
			} else if s.param.simAnnealFlag {
				pAcc := metropolis(s.pBest[idx].evalValue, solution.evalValue, s.param.t)
				if rand.Float64() < pAcc {
					s.pBest[idx].copy(solution)
				}
			}

			if solution.evalValue < s.pBestMem[idx].evalValue {
				s.pBestMem[idx].copy(solution)
			}
		}
		go worker(idx)
	}

	s.wg.Wait()
	s.updateBest()
	s.simAnneal()

	return err
}

func (s *swarm) updateBest() {
	for _, best := range s.pBestMem {
		if best.evalValue < s.gBest.evalValue {
			s.gBest = best
		}
	}
	if s.gBest.evalValue < s.gBestMem.evalValue {
		s.gBestMem.copy(s.gBest)
	}
}

func (s *swarm) getBestSolution() *Solution {
	return s.gBestMem
}

func (s *swarm) simAnneal() {
	s.param.t *= 0.999
}

func metropolis(eOld, eNew, T float64) float64 {
	return math.Exp((eOld - eNew) / T)
}
