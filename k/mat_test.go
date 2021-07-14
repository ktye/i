package mat

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"
)

func TestReal(t *testing.T) {
	A := [][]float64{
		[]float64{3, 5, -8, 12},
		[]float64{-2, 3, 3, 0},
		[]float64{7, -8, 2, 1},
	}

	qr := NewRQ(A)
	fmt.Println("qr A:", qr)
	r := qr.Solve([]float64{2, 3, 1, -5})
	fmt.Println("r:", r)
}
func TestMat(t *testing.T) {
	t.Skip()
	H := [][]complex128{
		[]complex128{1, 5i, 2, 4},
		[]complex128{-2i, 3, 3, -1},
		[]complex128{3, 2, 1, 1},
	}

	b := []complex128{1, 2, 3}
	for i := range b {
		b[i] *= cmplx.Rect(1, 30.0*math.Pi/180.0)
	}
	dot := func(A [][]complex128, j int, x []complex128) (r complex128) {
		for i := range x {
			r += A[i][j] * x[i]
		}
		return r
	}
	y := make([]complex128, 4)
	for i := range y {
		y[i] = dot(H, i, b)
	}

	qr := NewQR(H)
	fmt.Println("m/n", qr.M, qr.N)
	fmt.Println("Rdiag: ", p(qr.Rdiag))
	fmt.Println("H(")
	for _, h := range qr.H {
		fmt.Println(p(h))
	}
	fmt.Println(")")

	fmt.Println("y: ", p(y))
	r := qr.Solve(y)
	fmt.Println("r=A\\y: ", p(r))
}
