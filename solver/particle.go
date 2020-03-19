package solver

import (
	"libRPSO/vector"
	"math/rand"
)

type Bound struct {
	XUpper float64
	XLower float64
	VUpper float64
	VLower float64
}

type particle struct {
	solution *Solution
	velocity *Velocity
}

type InitParticleFunc func(bound *Bound, dim int) (*Solution, *Velocity)

func newParticle(param *PSOParam) (*particle, error) {
	if param == nil {
		panic("the param of particle can not be nil")
	}
	p := &particle{}
	p.init(param.dim)
	s, v := param.initFunc(param.bound, param.dim)
	p.velocity = v
	x, e, err := param.targetFunc(s.position, param.args)
	if err != nil {
		return nil, err
	}
	p.solution.update(x, e)
	return p, err
}

func (p *particle) init(dim int) {
	p.solution = NewSolution(dim)
	p.velocity = NewVelocity(dim)
}

func (p *particle) step(idx int, pBest []*Solution, gBest *Solution, param *PSOParam) error {
	p.solution.position = vector.Add(p.solution.position, p.velocity.v, 1, 1)
	p.updateVelocity(idx, pBest, gBest, param)

	pm := rand.Float64()
	if pm < param.pm {
		p.solution, _ = param.initFunc(param.bound, param.dim)
	}

	x, e, err := param.targetFunc(p.solution.position, param.args)
	if err != nil {
		return err
	}
	p.solution.update(x, e)

	return nil
}

func (p *particle) updateVelocity(idx int, pBest []*Solution, gBest *Solution, param *PSOParam) {
	// v = w*v
	p.velocity.v = vector.Scale(p.velocity.v, param.w)
	// v += c1*r1*(pBest[i] - pos)
	p.velocity.learn(pBest[idx].position, p.solution.position, param.c1)
	// v += c2*r2*(gBest - pos)
	p.velocity.learn(gBest.position, p.solution.position, param.c2)

	// v += c3*r3*(pBest[r] - pos)
	pr := rand.Float64()
	if pr < param.pr {
		rIdx := rand.Intn(len(pBest))
		p.velocity.learn(pBest[rIdx].position, p.solution.position, param.c3)
	}
}

func DefaultInitParticle(bound *Bound, dim int) (*Solution, *Velocity) {
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
