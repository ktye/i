package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"text/tabwriter"
)

func main() {
	A := [][]float64{
		{1, -2, 3},
		{5, 3, 2},
		{2, 3, 1},
		{4, -1, 1},
	}
	fmt.Println("A", mat(A))
	d, _ := NewReal(A)
	fmt.Println("H", mat(d.H))
	fmt.Println("R", d.Rdiag)

	b := []float64{1, 2, 3, 4}
	c, _ := d.QMul(b)
	fmt.Println("QTb", c)
	x, _ := d.RSolve(c)
	fmt.Println("x", x)
}
func mat(A [][]float64) string {
	var b bytes.Buffer
	n := 0
	if len(A) > 0 {
		n = len(A[0])
	}
	fmt.Fprintf(&b, "(%d,%d)\n", len(A), n)
	w := tabwriter.NewWriter(&b, 2, 8, 2, ' ', 0)
	for _, v := range A {
		for _, u := range v {
			fmt.Fprintf(w, "\t%.4f", u)
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return string(b.Bytes())
}

type RQ struct {
	H     [][]float64
	Rdiag []float64
	m, n  int
}

func NewReal(A [][]float64) (RQ, error) {
	m := len(A)    // Number of rows.
	n := len(A[0]) // Number of columns.
	if m < n {
		return RQ{}, errors.New("qr: matrix is underdetermined")
	}
	H := make([][]float64, n)
	Rdiag := make([]float64, n)
	for i := 0; i < n; i++ {
		H[i] = make([]float64, m)
		for k := 0; k < m; k++ {
			H[i][k] = A[k][i]
		}
	}
	for j := 0; j < n; j++ {
		fmt.Printf("H(%d) = %s\n", j, mat(H))
		s := RealNorm(H[j][j:])
		//fmt.Println("norm", s)
		if s == 0 {
			return RQ{}, errors.New("matrix contains zero-columns")
		}
		if H[j][j] > 0 {
			Rdiag[j] = -s
		} else {
			Rdiag[j] = s
		}
		//fmt.Printf("R[%d] = %v\n", j, Rdiag[j])
		f := 1.0 / math.Sqrt(s*(s+math.Abs(H[j][j])))
		//fmt.Println("f", f)
		H[j][j] -= Rdiag[j]
		//fmt.Println("Qii", H[j][j])
		for k := j; k < m; k++ {
			H[j][k] *= f
		}
		//fmt.Println("H", mat(H))
		//fmt.Println("Qrow", H[j])
		for i := j + 1; i < n; i++ {
			var sum float64
			for k := j; k < m; k++ {
				//fmt.Println("sum[%d]=%v\n", k, H[j][k]*H[i][k])
				sum += H[j][k] * H[i][k]
			}
			//fmt.Printf("sum (j+1)[%d] = %v\n", j+1, sum)
			for k := j; k < m; k++ {
				//fmt.Printf("H[%d][%d] -= %v\n", j, k, H[j][k]*sum)
				H[i][k] -= H[j][k] * sum
			}
		}
	}
	return RQ{
		H:     H,
		Rdiag: Rdiag,
		m:     m,
		n:     n,
	}, nil
}
func (D RQ) Solve(b []float64) ([]float64, error) {
	if len(b) != D.m {
		return nil, errors.New("qr: wrong input dimension for QR.Solve.")
	}
	if QTx, err := D.QMul(b); err != nil {
		return nil, err
	} else {
		return D.RSolve(QTx)
	}
}
func (D RQ) QMul(x []float64) ([]float64, error) {
	if len(x) != D.m {
		return nil, errors.New("qr: input vector lengths mismatch for QMul.")
	}
	y := make([]float64, D.m)
	for i := 0; i < D.m; i++ {
		y[i] = x[i]
	}
	for j := 0; j < D.n; j++ {
		var sum float64
		for k := j; k < D.m; k++ {
			sum += D.H[j][k] * y[k]
		}
		for k := j; k < D.m; k++ {
			y[k] -= D.H[j][k] * sum
		}
		//fmt.Printf("Qmul[%d] = %v\n", j, y)
	}
	return y, nil
}
func (D RQ) RSolve(b []float64) ([]float64, error) {
	if len(b) != D.m {
		return nil, errors.New("qr: input vector lengths mismatch for RSolve.")
	}
	x := make([]float64, D.m)
	for i := 0; i < D.m; i++ {
		x[i] = b[i]
	}
	for i := D.n - 1; i >= 0; i-- {
		s := 0.0
		for j := i + 1; j < D.n; j++ {
			fmt.Println("col ji", i+j*D.m)
			s += D.H[j][i] * x[j]
		}
		x[i] -= s
		fmt.Println("x", i, x)
		x[i] /= D.Rdiag[i]
		fmt.Println("X", i, x)
		//fmt.Printf("Rsolve[%d] = %v (D=%v)\n", i, x, D.Rdiag[i])
	}
	return x[0:D.n], nil
}
func RealNorm(v []float64) (r float64) {
	s := 0.0
	for _, x := range v {
		if x != 0 {
			x = math.Abs(x)
			if s < x {
				t := s / x
				r = 1 + r*t*t
				s = x
			} else {
				t := x / s
				r += t * t
			}
		}
	}
	return s * math.Sqrt(r)
}
