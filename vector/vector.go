package vector

import (
	"fmt"
	"math"
)

func CheckEqualLen(xs, ys []float64) {
	if len(xs) != len(ys) {
		panic(fmt.Sprintf("Lengths not equal: %d %d", len(xs), len(ys)))
	}
}

func DotProd(xs, ys []float64) (sum float64) {
	CheckEqualLen(xs, ys)
	for idx, x := range xs {
		sum += x * ys[idx]
	}
	return sum
}

func DocInPlace(xs, ys []float64) {
	for idx, x := range xs {
		xs[idx] = x * ys[idx]
	}
}

func Add(xs, ys []float64, alpha, beta float64) []float64 {
	CheckEqualLen(xs, ys)
	result := make([]float64, len(xs), len(xs))
	for idx, val := range xs {
		result[idx] = alpha*val + beta*ys[idx]
	}
	return result
}

func AddInPlace(accum, xs []float64, alpha float64) {
	CheckEqualLen(accum, xs)
	for idx, val := range xs {
		accum[idx] += alpha * val
	}
}

func Scale(xs []float64, alpha float64) []float64 {
	c := make([]float64, len(xs))
	for idx := range xs {
		c[idx] = alpha * xs[idx]
	}
	return c
}

func ScaleInPlace(xs []float64, alpha float64) {
	for idx, x := range xs {
		xs[idx] = alpha * x
	}
}

func lNorm(xs []float64, f func(x float64) float64) (norm float64) {
	for _, val := range xs {
		norm += f(val)
	}
	return
}

func L2(xs []float64) (magnitude float64) {
	return lNorm(xs, func(x float64) float64 {
		return x * x
	})
}

func L1(xs []float64) (magnitude float64) {
	return lNorm(xs, func(x float64) float64 {
		return math.Abs(x)
	})
}
