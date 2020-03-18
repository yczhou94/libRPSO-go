package libRPSO

import (
	"log"
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

type PSOParam struct {
	W        float64
	C1       float64
	C2       float64
	C3       float64
	Pr       float64
	Pm       float64
	T        float64
	Bound    *Bound
	Dim      int
	NProc    int
	PopSize  int
	InitFunc InitParticleFunc
	Target   TargetFunc
	Args     interface{}
}

func NewPSOParam() *PSOParam {
	return &PSOParam{
		W:  0.723,
		C1: 1.4454,
		C2: 1.4454,
		C3: 0.72,
		Pr: 0.1,
		Pm: 0.01,
		T:  100,
		Bound: &Bound{
			XUpper: 10,
			XLower: -10,
			VUpper: 5,
			VLower: -5,
		},
		NProc: 1,
	}
}

func newSwarm(param *PSOParam) (*swarm, error) {
	s := &swarm{
		popSize:   param.PopSize,
		particles: make([]*particle, param.PopSize),
		pBest:     make([]*Solution, param.PopSize),
		gBest:     NewSolution(param.Dim),
		pBestMem:  make([]*Solution, param.PopSize),
		gBestMem:  NewSolution(param.Dim),
		param:     param,
	}

	for i := 0; i < param.PopSize; i++ {
		p, err := NewParticle(s.param)
		if err != nil {
			return nil, err
		}
		s.particles[i] = p
		s.pBest[i] = NewSolution(param.Dim)
		s.pBestMem[i] = NewSolution(param.Dim)
		s.pBest[i].copy(p.solution)
	}

	copy(s.pBestMem, s.pBest)
	s.updateBest()

	return s, nil
}

func (s *swarm) step() error {
	var err error
	s.wg.Add(s.param.PopSize)
	queue := make(chan struct{}, s.param.NProc)
	for idx, p := range s.particles {
		queue <- struct{}{}
		worker := func(idx int) {
			defer func() {
				<-queue
				s.wg.Done()
			}()
			log.Printf("particle %d running \n", idx)
			e := p.step(idx, s.pBest, s.gBest, s.param)
			if e != nil {
				err = e
				return
			}
			if p.solution.evalValue < s.pBest[idx].evalValue {
				s.pBest[idx].copy(p.solution)
			} else {
				pAcc := metropolis(s.pBest[idx].evalValue, p.solution.evalValue, s.param.T)
				if pAcc < rand.Float64() {
					s.pBest[idx].copy(p.solution)
				}
			}

			if p.solution.evalValue < s.pBestMem[idx].evalValue {
				s.pBestMem[idx].copy(p.solution)
			}
		}
		go worker(idx)
	}

	s.wg.Wait()
	s.updateBest()

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

func metropolis(eOld, eNew, T float64) float64 {
	return math.Exp((eOld - eNew) / T)
}
