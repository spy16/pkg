package matrix_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spy16/pkg/matrix"
)

func TestMatrix_Clone(t *testing.T) {
	A := matrix.FromSlice(2, 3, []float64{
		1, 2, 3,
		3, 4, 5,
	})
	cloned := A.Clone()

	if !reflect.DeepEqual(cloned, A) {
		t.Errorf("cloned copy was '%#v', but expected '%#v'", cloned, A)
	}
}

func TestFromSlice(t *testing.T) {
	A := matrix.FromSlice(2, 3, []float64{
		1, 2, 3,
		3, 4, 5,
	})

	r, c := A.Dims()
	if r != 2 || c != 3 {
		t.Errorf("expecting dimesion to be 2x3, got %dx%d", r, c)
	}

	e1 := A.Elem(0, 1)
	if e1 != 2 {
		t.Errorf("expecting element at (%d,%d) to be %f, got %f", 0, 1, 2., e1)
	}
}

func TestProduct(suite *testing.T) {
	suite.Parallel()

	cases := []struct {
		A, B, Res *matrix.Matrix
	}{
		{
			A: matrix.FromSlice(2, 2, []float64{
				1, 2,
				5, 6,
			}),
			B: matrix.FromSlice(2, 1, []float64{
				5,
				6,
			}),
			Res: matrix.FromSlice(2, 1, []float64{
				17,
				61,
			}),
		},
		{
			A: matrix.FromSlice(2, 2, []float64{
				1, 2,
				5, 6,
			}),
			B: matrix.FromSlice(2, 2, []float64{
				1, 2,
				5, 6,
			}),
			Res: matrix.FromSlice(2, 2, []float64{
				11, 14,
				35, 46,
			}),
		},
	}

	for id, cs := range cases {
		suite.Run(fmt.Sprintf("#%d", id), func(t *testing.T) {
			res := matrix.Product(cs.A, cs.B)

			if !reflect.DeepEqual(res, cs.Res) {
				t.Errorf("expected %#v, got %#v", cs.Res, res)
			}
		})
	}
}
