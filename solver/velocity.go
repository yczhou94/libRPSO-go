package solver

import (
	"libRPSO/vector"
	"math/rand"
)

type Velocity struct {
	v []float64
}

func NewVelocity(dim int) *Velocity {
	return &Velocity{
		v: make([]float64, dim),
	}
}

func (v *Velocity) learn(px, py []float64, c float64) {
	tmp := vector.Add(px, py, 1, -1)
	r := rand.Float64()
	tmp = vector.Scale(tmp, c*r)
	v.v = vector.Add(v.v, tmp, 1, 1)
}
