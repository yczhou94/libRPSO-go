package solver

import (
	"math"
	"math/rand"
	"sync"
)

type swarm struct {
	popSize   int
	particles []*particle
	pBest     []*Solution
	gBest     *Solution
	pBestMem  []*Solution
	gBestMem  *Solution
	param     *PSOParam
	wg        sync.WaitGroup
}

func newSwarm(param *PSOParam) (*swarm, error) {
	s := &swarm{
		popSize:   param.popSize,
		particles: make([]*particle, param.popSize),
		pBest:     make([]*Solution, param.popSize),
		gBest:     NewSolution(param.dim),
		pBestMem:  make([]*Solution, param.popSize),
		gBestMem:  NewSolution(param.dim),
		param:     param,
	}

	for i := 0; i < param.popSize; i++ {
		p, err := newParticle(s.param)
		if err != nil {
			return nil, err
		}
		s.particles[i] = p
		s.pBest[i] = NewSolution(param.dim)
		s.pBestMem[i] = NewSolution(param.dim)
		s.pBest[i].copy(p.solution)
	}

	copy(s.pBestMem, s.pBest)
	s.updateBest()

	return s, nil
}

func (s *swarm) step(step int) error {
	var err error
	s.wg.Add(s.param.popSize)
	queue := make(chan struct{}, s.param.nProc)
	for idx, p := range s.particles {
		queue <- struct{}{}
		worker := func(idx int) {
			defer func() {
				<-queue
				s.wg.Done()
			}()
			e := p.step(idx, s.pBest, s.gBest, s.param)
			if e != nil {
				err = e
				return
			}
			if p.solution.evalValue < s.pBest[idx].evalValue {
				s.pBest[idx].copy(p.solution)
			} else {
				pAcc := metropolis(s.pBest[idx].evalValue, p.solution.evalValue, s.param.t)
				if rand.Float64() < pAcc {
					s.pBest[idx].copy(p.solution)
				}
			}

			if p.solution.evalValue < s.pBestMem[idx].evalValue && s.param.simAnnealFlag {
				s.pBestMem[idx].copy(p.solution)
			}
		}
		go worker(idx)
	}

	s.wg.Wait()
	s.updateBest()
	s.simAnneal(step)

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

func (s *swarm) simAnneal(step int) {
	s.param.t /= 1 + math.Log(float64(step+1+1))
}

func metropolis(eOld, eNew, T float64) float64 {
	return math.Exp((eOld - eNew) / T)
}
