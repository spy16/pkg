package ffnet

import (
	"math"

	"github.com/spy16/pkg/ffnet/mat"
)

// Differentiable is used to contain an activation function and its
// derivative.
type Differentiable struct {
	F      func(x mat.Matrix) mat.Matrix
	FPrime func(z mat.Matrix) mat.Matrix
}

// LossFunc is used to contain a loss function and its derivative.
type LossFunc struct {
	F      func(y, yHat mat.Matrix) mat.Matrix
	FPrime func(y, yHat mat.Matrix) mat.Matrix
}

// Sigmoid is a logistic activator in the special case of a = 1.
func Sigmoid() Differentiable {
	return Differentiable{
		F: func(x mat.Matrix) mat.Matrix {
			// result = 1.0 / (1.0 + exp(-x))
			sigmoid := func(_, _ int, val float64) float64 {
				return 1.0 / (1.0 + math.Exp(-val))
			}
			return x.Apply(sigmoid)
		},
		FPrime: func(z mat.Matrix) mat.Matrix {
			// res = 1-z
			res := z.Clone().Apply(func(_, _ int, val float64) float64 {
				return 1 - val
			})
			// return z * (1-z)
			return mat.Mul(z, res)
		},
	}
}

// ReLU implements the Rectified Linear Unit function.
func ReLU() Differentiable {
	return Differentiable{
		F: func(x mat.Matrix) mat.Matrix {
			// result = max(0, x)
			return x.Apply(func(_, _ int, val float64) float64 {
				return math.Max(0, val)
			})
		},
		FPrime: func(z mat.Matrix) mat.Matrix {
			// result = 1 if x > 0 else 0
			return z.Apply(func(_, _ int, val float64) float64 {
				if val > 0 {
					return 1
				}
				return 0
			})
		},
	}
}

// SquaredError represents squared difference between desired and actual
// prediction as the error.
func SquaredError() LossFunc {
	return LossFunc{
		F: func(y, yHat mat.Matrix) mat.Matrix {
			// return (1/2)*(y-yHat).^2
			return mat.Sub(y, yHat).Apply(mat.Square).Scale(1.0 / 2.0)
		},
		FPrime: func(y, yHat mat.Matrix) mat.Matrix {
			// return (yHat-y)
			return mat.Sub(yHat, y)
		},
	}
}
