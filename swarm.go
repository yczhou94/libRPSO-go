package libRPSO

import (
	"math"
	"math/rand"
)

type swarm struct {
	popSize   int
	particles []*particle
	pBest     []*Solution
	gBest     *Solution
	pBestMem  []*Solution
	gBestMem  *Solution
	param     *PSOParam
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
	return &PSOParam{}
}

func NewSwarm(param *PSOParam) (*swarm, error) {
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
	for idx, p := range s.particles {
		err := p.step(idx, s.pBest, s.gBest, s.param)
		if err != nil {
			return err
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

	s.updateBest()

	return nil
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
