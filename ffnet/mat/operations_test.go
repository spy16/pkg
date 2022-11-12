package mat_test

import (
	"testing"

	"github.com/spy16/pkg/ffnet/mat"
)

func TestMatrix_Dot(t *testing.T) {
	m1 := mat.From(3, 2, []float64{
		1, 2,
		3, 4,
		5, 6,
	}...)

	m2 := mat.From(2, 3, []float64{
		1, 2, 3,
		4, 5, 6,
	}...)

	expected := mat.From(3, 3, []float64{
		9, 12, 15,
		19, 26, 33,
		29, 40, 51,
	}...)
	prod := mat.Dot(m1, m2)
	if !expected.Equals(prod) {
		t.Errorf("Matrix.Mul() want=%+v\ngot=%+v", expected, prod)
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	t.Run("SimpleAdd", func(t *testing.T) {
		m := mat.From(3, 2, []float64{
			1, 2,
			3, 4,
			5, 6,
		}...)

		want := mat.From(3, 2, []float64{
			2, 4,
			6, 8,
			10, 12,
		}...)
		got := mat.Add(m, m)

		if !want.Equals(got) {
			t.Errorf("Add() want=%v, got=%v", want, got)
		}
	})

	t.Run("ScalarAdd", func(t *testing.T) {
		m := mat.From(3, 2, []float64{
			1, 2,
			3, 4,
			5, 6,
		}...)
		s := mat.From(1, 1, 1)

		want := mat.From(3, 2, []float64{
			2, 3,
			4, 5,
			6, 7,
		}...)
		got := mat.Add(m, s)

		if !want.Equals(got) {
			t.Errorf("Add() want=%v, got=%v", want, got)
		}
	})

	t.Run("BroadcastRow", func(t *testing.T) {
		m := mat.From(3, 2, []float64{
			1, 2,
			3, 4,
			5, 6,
		}...)
		s := mat.From(1, 2, []float64{2, 3}...)

		want := mat.From(3, 2, []float64{
			3, 5,
			5, 7,
			7, 9,
		}...)
		got := mat.Add(m, s)

		if !want.Equals(got) {
			t.Errorf("Add() want=%v, got=%v", want, got)
		}
	})

	t.Run("BroadcastCol", func(t *testing.T) {
		m := mat.From(3, 2, []float64{
			1, 2,
			3, 4,
			5, 6,
		}...)
		s := mat.From(3, 1, []float64{1, 2, 3}...)

		want := mat.From(3, 2, []float64{
			2, 3,
			5, 6,
			8, 9,
		}...)
		got := mat.Add(m, s)

		if !want.Equals(got) {
			t.Errorf("Add() want=%v, got=%v", want, got)
		}
	})
}
