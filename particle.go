package libRPSO

import (
	"libRPSO/vector"
)

type Particle struct {
	position  []float64
	velocity  []float64
	evalValue float64
}

func NewParticle(pos []float64, v []float64, e float64) *Particle {
	if pos == nil {
		panic("position of particle can not be nil")
	}
	if v == nil {
		panic("velocity of particle can not be nil")
	}
	vector.CheckEqualLen(pos, v)
	return &Particle{
		position:  pos,
		velocity:  v,
		evalValue: e,
	}
}

func (p *Particle) Step(f TargetFunc, args ...interface{}) {

}
