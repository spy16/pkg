package mat

import (
	"fmt"
	"strings"
)

// New returns a new zero-valued matrix of given dimensions. Use Apply() to
// fill the matrix with values.
func New(rows, cols int) Matrix {
	return Matrix{
		rows: rows,
		cols: cols,
		vals: make([]float64, rows*cols, rows*cols),
	}
}

// From creates a matrix of given dimensions with values set to vals.
func From(rows, cols int, vals ...float64) Matrix {
	if len(vals) != (rows * cols) {
		panic(fmt.Errorf("need exactly %d values for matrix of dimensions=(%d,%d), got %d",
			rows*cols, rows, cols, len(vals)))
	}
	return Matrix{
		rows: rows,
		cols: cols,
		vals: vals,
	}
}

// Matrix represents a dense matrix of float64 values.
type Matrix struct {
	rows, cols int
	vals       []float64
}

// Equals returns true if other has dimensions and same elements.
func (m Matrix) Equals(other Matrix) bool {
	if m.Dims() != other.Dims() {
		return false
	}
	for i := 0; i < len(m.vals); i++ {
		if m.vals[i] != other.vals[i] {
			return false
		}
	}
	return true
}

// Values returns all values in the matrix as a 1d vector/slice.
func (m Matrix) Values() []float64 { return append([]float64(nil), m.vals...) }

// Clone returns an exact clone of m.
func (m Matrix) Clone() Matrix {
	return Matrix{
		rows: m.rows,
		cols: m.cols,
		vals: append([]float64(nil), m.vals...),
	}
}

// Apply updates the matrix in-place using the mapper function given. Returns
// the same matrix pointer for convenience. If `fn` is nil, matrix is reset to
// zero values.
func (m Matrix) Apply(fn Mapper) Matrix {
	if fn == nil {
		m.vals = make([]float64, len(m.vals), len(m.vals))
		return m
	}

	for i := 0; i < len(m.vals); i++ {
		r, c := m.cell(i)
		m.vals[i] = fn(r, c, m.vals[i])
	}
	return m
}

// T returns the transpose of m.
func (m Matrix) T() Matrix {
	t := New(m.cols, m.rows) // swapped dimensions
	m.Apply(func(row, col int, val float64) float64 {
		t.vals[t.index(col, row)] = val
		return val
	})
	return t
}

// Elem returns the value of the given cell.
func (m Matrix) Elem(row, col int) float64 { return m.vals[m.index(row, col)] }

// Set sets the value of the given cell.
func (m Matrix) Set(row, col int, val float64) { m.vals[m.index(row, col)] = val }

// Dims returns the dimension of m.
func (m Matrix) Dims() [2]int { return [2]int{m.rows, m.cols} }

// Size returns the total number of elements in the matrix.
func (m Matrix) Size() int { return len(m.vals) }

// Scale can be used multiply all values in a matrix by a scalar factor.
func (m Matrix) Scale(factor float64) Matrix {
	return m.Apply(func(_, _ int, val float64) float64 {
		return val * factor
	})
}

func (m Matrix) String() string {
	var s strings.Builder
	for i := 0; i < m.rows; i++ {
		_, _ = fmt.Fprintf(&s, "%v", m.row(i))
		if i < m.rows-1 {
			s.WriteRune('\n')
		}
	}
	return s.String()
}

func (m Matrix) row(i int) []float64 {
	items := make([]float64, m.cols)
	for j := 0; j < m.cols; j++ {
		items[j] = m.vals[m.index(i, j)]
	}
	return items
}

func (m Matrix) index(row, col int) int {
	if row >= m.rows {
		panic(fmt.Errorf("row %d is out of bounds [rows=%d]", row, m.rows))
	} else if col >= m.cols {
		panic(fmt.Errorf("column %d is out of bounds [cols=%d]", col, m.cols))
	}
	return row*m.cols + col
}

func (m Matrix) cell(idx int) (row, col int) {
	if idx >= len(m.vals) {
		panic(fmt.Errorf("1D index %d is out-of-range for matrix with dimensions=%dx%d",
			idx, m.rows, m.cols))
	}
	col = idx % m.cols
	row = (idx - col) / m.cols
	return row, col
}
