package main

import (
	"bytes"
	"fmt"
	"math"
	"math/cmplx"
	"strings"
	"text/tabwriter"
)

func main() {
	A := [][]complex128{
		{1, -2i, 3},
		{5i, 3, 2},
		{2, 3, 1},
		{4, -1, 1},
	}
	fmt.Println("A", mat(A))
	d := New(A)
	fmt.Println("H", mat(d.H))
	fmt.Println("D", vec(d.Rdiag))

	b := []complex128{1, 2, 3, 4}
	/*
		c := d.QMul(b)
		fmt.Println("c", vec(c))
		x := d.RSolve(c)
		fmt.Println("x", vec(x))
	*/
	fmt.Println("x", vec(d.Solve(b)))
}
func mat(A [][]complex128) string {
	var b bytes.Buffer
	n := 0
	if len(A) > 0 {
		n = len(A[0])
	}
	fmt.Fprintf(&b, "(%d,%d)\n", len(A), n)
	w := tabwriter.NewWriter(&b, 2, 8, 2, ' ', 0)
	for _, v := range A {
		for _, u := range v {
			fmt.Fprintf(w, "\t%s", absang(u))
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return string(b.Bytes())
}
func vec(A []complex128) string {
	v := make([]string, len(A))
	for i, u := range A {
		v[i] = absang(u)
	}
	return "[" + strings.Join(v, " ") + "]"
}
func absang(z complex128) string {
	a := cmplx.Phase(z) / math.Pi * 180.0
	if a < 0 {
		a += 360.0
	}
	return fmt.Sprintf("%.4fa%.0f", cmplx.Abs(z), a)
}

type QR struct {
	H     [][]complex128
	Rdiag []complex128 // Missing diagonal of R.
	m, n  int          // number of rows and columns
}

// New calculates the QR Decomposition of a rectangular matrix.
func New(A [][]complex128) QR {
	m := len(A)    // Number of rows.
	n := len(A[0]) // Number of columns.
	if m < n {
		panic("qr: matrix is underdetermined")
	}

	H := make([][]complex128, n)
	Rdiag := make([]complex128, n)
	for i := 0; i < n; i++ {
		H[i] = make([]complex128, m)
		for k := 0; k < m; k++ {
			H[i][k] = A[k][i]
		}
	}
	for j := 0; j < n; j++ {
		s := VectorNorm(H[j][j:])
		if s == 0 {
			panic("matrix contains zero-columns")
		}

		Rdiag[j] = -complex(s, 0) * cmplx.Rect(1, cmplx.Phase(H[j][j])) // Diagonal element.
		f := complex(math.Sqrt(s*(s+cmplx.Abs(H[j][j]))), 0)
		H[j][j] -= Rdiag[j]

		for k := j; k < m; k++ {
			H[j][k] /= f
		}
		for i := j + 1; i < n; i++ {
			var sum complex128
			for k := j; k < m; k++ {
				sum += cmplx.Conj(H[j][k]) * H[i][k]
			}
			for k := j; k < m; k++ {
				H[i][k] -= H[j][k] * sum
			}
		}
	}
	return QR{
		H:     H,
		Rdiag: Rdiag,
		m:     m,
		n:     n,
	}
}
func (D QR) Solve(b []complex128) []complex128 {
	if len(b) != D.m {
		panic("qr: wrong input dimension for QR.Solve.")
	}
	return D.RSolve(D.QMul(b))
}
func (D QR) QMul(x []complex128) []complex128 {
	if len(x) != D.m {
		panic("qr: input vector lengths mismatch for QMul.")
	}
	y := make([]complex128, D.m)
	for i := 0; i < D.m; i++ {
		y[i] = x[i]
	}
	for j := 0; j < D.n; j++ {
		var sum complex128
		for k := j; k < D.m; k++ {
			sum += cmplx.Conj(D.H[j][k]) * y[k]
		}
		for k := j; k < D.m; k++ {
			y[k] -= D.H[j][k] * sum
		}
	}
	return y
}
func (D QR) RSolve(b []complex128) []complex128 {
	fmt.Println("solve b=", vec(b))
	if len(b) != D.m {
		panic("qr: input vector lengths mismatch for RSolve.")
	}
	x := make([]complex128, D.m)
	for i := 0; i < D.m; i++ {
		x[i] = b[i]
	}
	for i := D.n - 1; i >= 0; i-- {
		var s complex128
		//fmt.Printf("s(%d) = +/", i)
		for j := i + 1; j < D.n; j++ {
			//fmt.Printf(" (%s*%s)=%s", absang(D.H[j][i]), absang(x[j]), absang(D.H[j][i]*x[j]))
			s += D.H[j][i] * x[j]
		}
		//fmt.Printf("\ns = %s\n", absang(s))
		//fmt.Printf("x[%d]= %s -s = %s\n", i, absang(x[i]), absang(x[i]-s))
		x[i] -= s
		fmt.Printf("x[%d]=%s / %s = %s\n", i, absang(x[i]), absang(D.Rdiag[i]), absang(x[i]/D.Rdiag[i]))
		x[i] /= D.Rdiag[i]
	}
	return x[0:D.n]
}
func VectorNorm(x []complex128) (norm float64) {
	for i := 0; i < len(x); i++ {
		norm = math.Hypot(norm, cmplx.Abs(x[i]))
	}
	return
}
