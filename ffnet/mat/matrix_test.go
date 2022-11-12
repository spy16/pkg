package mat_test

import (
	"testing"

	"github.com/spy16/pkg/ffnet/mat"
)

func TestNew(t *testing.T) {
	m := mat.New(3, 3).Apply(func(row, col int, val float64) float64 { return 1.0 })

	if m.Size() != 9 {
		t.Errorf("New() expected size to be 9, not %d", m.Size())
	}

	if m.Dims() != [2]int{3, 3} {
		t.Errorf("New() expected dimensions to be [3 3], not %v", m.Dims())
	}

	m.Apply(func(row, col int, val float64) float64 {
		if val != 1. {
			t.Errorf("New() expected value of cell (%d, %d) to be 1.0, not %f", row, col, val)
		}
		return val
	})
}

func TestMatrix_Scale(t *testing.T) {
	m := mat.From(3, 3, []float64{
		1, 0, 2,
		0, 1, 0,
		0, 0, 1,
	}...)
	want := mat.From(3, 3, []float64{
		-1, 0, -2,
		0, -1, 0,
		0, 0, -1,
	}...)

	scaled := m.Clone().Scale(-1)
	if !scaled.Equals(want) {
		t.Errorf("Matrix.Scale() want=%v, got%v", want, scaled)
	}
}

func TestMatrix_Equals(t *testing.T) {
	t.Parallel()

	t.Run("Same", func(t *testing.T) {
		m := mat.From(3, 3, []float64{
			1, 0, 2,
			0, 1, 0,
			0, 0, 1,
		}...)
		if !m.Equals(m) {
			t.Errorf("Matrix.Equals() want true, got false")
		}
	})

	t.Run("DifferentVals", func(t *testing.T) {
		m := mat.From(3, 3, []float64{
			1, 0, 2,
			0, 1, 0,
			0, 0, 1,
		}...)
		if m.Equals(m.Clone().Apply(mat.Ones)) {
			t.Errorf("Matrix.Equals() want false, got true")
		}
	})

	t.Run("DifferentDims", func(t *testing.T) {
		m := mat.From(3, 2, []float64{
			1, 0,
			0, 1,
			0, 0,
		}...)
		if m.Equals(m.T()) {
			t.Errorf("Matrix.Equals() want false, got true")
		}
	})
}

func TestMatrix_T(t *testing.T) {
	t.Parallel()

	t.Run("Square", func(t *testing.T) {
		m := mat.From(3, 3, []float64{
			1, 0, 2,
			0, 1, 0,
			0, 0, 1,
		}...)
		mT := m.T()
		expected := mat.From(3, 3, []float64{
			1, 0, 0,
			0, 1, 0,
			2, 0, 1,
		}...)
		if !mT.Equals(expected) {
			t.Errorf("Matrix.T() want=%v, got=%v", expected, mT)
		}
	})

	t.Run("NonSquare", func(t *testing.T) {
		m := mat.From(3, 2, []float64{
			1, 2,
			3, 4,
			5, 6,
		}...)
		mT := m.T()
		expected := mat.From(2, 3, []float64{
			1, 3, 5,
			2, 4, 6,
		}...)
		if !mT.Equals(expected) {
			t.Errorf("Matrix.T() want=%v, got=%v", expected, mT)
		}
	})
}

func TestMatrix_Clone(t *testing.T) {
	m := mat.From(3, 2, []float64{
		1, 2,
		3, 4,
		5, 6,
	}...)

	c := m.Clone()

	if !m.Equals(c) {
		t.Errorf("Matrix.Clone() want=%+v\ngot=%+v", m, c)
	}
}

func TestMatrix_Set(t *testing.T) {
	m := mat.From(3, 2, []float64{
		1, 2,
		3, 4,
		5, 6,
	}...)

	m.Set(0, 0, 100)
	if m.Elem(0, 0) != 100 {
		t.Errorf("Matrix.Elem(0, 0) want=%+v\ngot=%+v", 100, m.Elem(0, 0))
	}
}

func BenchmarkNew_Random(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mat.New(100, 100).Apply(mat.Random)
	}
}
