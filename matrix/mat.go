package matrix

import (
	"fmt"
	"math/rand"
	"time"
)

// FromSlice creates a matrix of given dimensions from the given slice.
func FromSlice(rows, cols int, arr []float64) *Matrix {
	L := rows * cols
	if L != len(arr) {
		panic(fmt.Errorf("slice must be of length %d for a matrix of dimension %dx%d, but got %d", L, rows, cols, len(arr)))
	}

	m := &Matrix{
		rows: rows,
		cols: cols,
		data: map[int]float64{},
	}
	for i, itm := range arr {
		m.data[i] = itm
	}
	return m
}

// New initializes a sparse matrix with dimension (rows,cols) and
// given zero value.
func New(rows, cols int, zeroVal float64) *Matrix {
	return &Matrix{
		rows:    rows,
		cols:    cols,
		zeroVal: zeroVal,
		data:    map[int]float64{},
	}
}

// Random initializes a random-value filled matrix with given dimensions.
// zeroVal can be used to control zero value of the sparse matrix. Use 0
// for zeroVal by default.
func Random(rows, cols int, zeroVal float64, randSrc rand.Source) *Matrix {
	m := New(rows, cols, zeroVal)

	if randSrc == nil {
		randSrc = rand.NewSource(time.Now().UnixNano())
	}

	r := rand.New(randSrc)
	for i := 0; i < rows*cols; i++ {
		v := r.Float64()
		if v != zeroVal {
			m.data[i] = v
		}
	}

	return m
}

// Zeros initializes a zero-filled matrix with given dimensions.
func Zeros(rows, cols int) *Matrix {
	return &Matrix{
		rows:    rows,
		cols:    cols,
		data:    map[int]float64{},
		zeroVal: 0,
	}
}

// Product calculates the matrix product of the given matrices.
func Product(A, B *Matrix) *Matrix {
	aR, aC := A.Dims()
	bR, bC := B.Dims()

	if aC != bR {
		panic(fmt.Errorf("multiplication not possible, inner dimensions mismatch"))
	}

	pr := Zeros(aR, bC)

	for i := 0; i < aR; i++ {
		for j := 0; j < bC; j++ {
			p := 0.0
			for k := 0; k < aC; k++ {
				p += A.Elem(i, k) * B.Elem(k, j)
			}
			pr.Set(i, j, p)
		}
	}

	return pr
}

// Reduce provides a functional way of reducing a function over the whole matrix.
// For instance: sum can be implemented as:
//   m.Reduce(0, func(x, cum float64) float64 {return x+cum})
func Reduce(m *Matrix, zero float64, apply func(x, acc float64) float64) float64 {
	acc := zero
	for _, x := range m.data {
		acc = apply(x, acc)
	}
	return acc
}

// Apply clones the matrix and returns the clone after applying the given function
// to each element of the matrix.
func Apply(m *Matrix, f func(float64) float64) *Matrix {
	clone := m.Clone()
	for i := 0; i < m.Size(); i++ {
		v := f(clone.at(i))
		if v != clone.zeroVal {
			clone.data[i] = v
		}
	}

	return clone
}

// Sum returns the sum of all elements in the given matrix.
func Sum(m *Matrix) float64 {
	return Reduce(m, 0, func(x, acc float64) float64 {
		return acc + x
	})
}

// Dot calculates the dot product of given matrices.
func Dot(a, b *Matrix) *Matrix {
	return dotApply(a, b, func(x, y float64) float64 { return x * y })
}

// Plus returns the sum of the two matrices.
func Plus(a, b *Matrix) *Matrix {
	return dotApply(a, b, func(x, y float64) float64 { return x + y })
}

// Minus returns the difference of the two matrices.
func Minus(a, b *Matrix) *Matrix {
	return dotApply(a, b, func(x, y float64) float64 { return x - y })
}

// Matrix represents a sparse matrix with dimensions (rows x cols).
// Matrix is not safe for concurrent access.
type Matrix struct {
	data       map[int]float64
	rows, cols int
	zeroVal    float64
}

// Dims returns the dimensions of the matrix.
func (m *Matrix) Dims() (rows, cols int) {
	return m.rows, m.cols
}

// Size returns the underlying data array size.
func (m *Matrix) Size() int {
	return m.rows * m.cols
}

// Elem returns the element at the given row and col.
func (m *Matrix) Elem(row, col int) float64 {
	return m.at(m.index(row, col))
}

// Set sets the value of the cell at (row, col).
func (m *Matrix) Set(row, col int, val float64) {
	if val == m.zeroVal {
		delete(m.data, m.index(row, col))
	} else {
		m.data[m.index(row, col)] = val
	}
}

// Unset sets the cell value to configured zero value and returns
// the current value.
func (m *Matrix) Unset(row, col int) float64 {
	v := m.Elem(row, col)
	m.data[m.index(row, col)] = m.zeroVal
	return v
}

// Clone returns a copy of the matrix.
func (m *Matrix) Clone() *Matrix {
	cloned := &Matrix{
		rows: m.rows,
		cols: m.cols,
		data: map[int]float64{},
	}

	for k, v := range m.data {
		cloned.data[k] = v
	}

	return cloned
}

func (m *Matrix) at(i int) float64 {
	val, ok := m.data[i]
	if !ok {
		return m.zeroVal
	}
	return val
}

func (m *Matrix) index(row, col int) int {
	if row < 0 || row >= m.rows {
		panic(fmt.Errorf("row index '%d' out of bounds, must be in range (0,%d]", row, m.rows))
	}

	if col < 0 || col >= m.cols {
		panic(fmt.Errorf("column index '%d' out of bounds, must be in range (0,%d]", col, m.cols))
	}

	return row*m.cols + col
}

func dotApply(a, b *Matrix, f func(x, y float64) float64) *Matrix {
	if a.cols != b.cols || a.rows != b.rows {
		panic(fmt.Sprintf("can't compute dot application of matrices with dimensions %dx%d and %dx%d",
			a.rows, a.cols, b.rows, b.cols))
	}

	c := Zeros(a.rows, a.cols)
	for i := 0; i < a.rows; i++ {
		for j := 0; j < b.cols; j++ {
			c.Set(i, j, f(a.Elem(i, j), b.Elem(i, j)))
		}
	}
	return c
}
