package mat

import (
	"fmt"
	"testing"
)

func TestMat(t *testing.T) {
	H := [][]complex128{
		[]complex128{1, 5i, 2, 4},
		[]complex128{-2i, 3, 3, -1},
		[]complex128{3, 2, 1, 1},
	}
	qr := NewQR(H)
	fmt.Println("m/n", qr.M, qr.N)
	fmt.Println("Rdiag: ", p(qr.Rdiag))
	fmt.Println("H(")
	for _, h := range qr.H {
		p(h)
	}
	fmt.Println(")")
}
