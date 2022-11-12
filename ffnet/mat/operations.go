package mat

import "fmt"

// Dot performs dot product of a & b and stores the result into m. Panics
// if inner dimensions are not same.
func Dot(a, b Matrix) Matrix {
	if a.cols != b.rows {
		panic(fmt.Errorf("inner dimensions mismatch: %v & %v", a.Dims(), b.Dims()))
	}

	return New(a.rows, b.cols).Apply(func(i, j int, val float64) float64 {
		var sum float64
		for k := 0; k < a.cols; k++ {
			sum += a.vals[a.index(i, k)] * b.vals[b.index(k, j)]
		}
		return sum
	})
}

// Add computes a + b with broadcasting support and returns the result. Panics
// if can't be broadcast.
func Add(a, b Matrix) Matrix { return broadcastOp('+', a, b) }

// Sub computes a - b with broadcasting support and returns the result. Panics
// if can't be broadcast.
func Sub(a, b Matrix) Matrix { return broadcastOp('-', a, b) }

// Mul computes a * b with broadcasting support and returns the result. Panics
// if can't be broadcast.
func Mul(a, b Matrix) Matrix { return broadcastOp('*', a, b) }

// Div computes a / b with broadcasting support and returns the result. Panics
// if can't be broadcast.
func Div(a, b Matrix) Matrix { return broadcastOp('/', a, b) }

func broadcastOp(op rune, a, b Matrix) Matrix {
	hiDims := higherDims(a, b)
	a = broadcastToHi(a, hiDims)
	b = broadcastToHi(b, hiDims)

	m := New(a.rows, a.cols)
	for i := 0; i < len(a.vals); i++ {
		switch op {
		case '+':
			m.vals[i] = a.vals[i] + b.vals[i]

		case '-':
			m.vals[i] = a.vals[i] - b.vals[i]

		case '*':
			m.vals[i] = a.vals[i] * b.vals[i]

		case '/':
			m.vals[i] = a.vals[i] / b.vals[i]

		default:
			panic(fmt.Errorf("'%c' is not valid op", op))
		}
	}
	return m
}

func broadcastToHi(lo Matrix, hiDims [2]int) Matrix {
	loDims := lo.Dims()

	switch {
	case hiDims == loDims:
		// already same
		return lo

	case loDims == [2]int{1, 1}:
		// scalar value
		return New(hiDims[0], hiDims[1]).Apply(Value(lo.vals[0]))

	case loDims[0] == hiDims[0] && loDims[1] == 1:
		// column vector
		return New(hiDims[0], hiDims[1]).Apply(func(row, col int, val float64) float64 {
			return lo.vals[lo.index(row, 0)]
		})

	case loDims[1] == hiDims[1] && loDims[0] == 1:
		// row vector
		return New(hiDims[0], hiDims[1]).Apply(func(row, col int, val float64) float64 {
			return lo.vals[lo.index(0, col)]
		})

	default:
		panic(fmt.Errorf("cannot broadcast %v to %v", loDims, hiDims))
	}
}

func higherDims(a Matrix, b Matrix) [2]int {
	if len(a.vals) > len(b.vals) {
		return a.Dims()
	}
	return b.Dims()
}
