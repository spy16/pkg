package mat

import "math/rand"

// Mapper implementation can be used with New() or Apply() to perform initialize
// or to update a matrix.
type Mapper func(row, col int, val float64) float64

// Random can be used with New() or Apply() to fill a matrix with random values
// provided by rand.Float64().
func Random(_, _ int, _ float64) float64 { return rand.Float64() }

// Ones can be used with New() or Apply() to fill a matrix with ones.
func Ones(_, _ int, _ float64) float64 { return 1 }

// Value returns a mapper that updates the matrix with a constant value.
func Value(v float64) Mapper {
	return func(_, _ int, _ float64) float64 { return v }
}

// Square can be used with Apply() to square all values in a matrix.
func Square(_, _ int, v float64) float64 {
	return v * v
}
