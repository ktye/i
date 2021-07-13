package mat

import (
	"fmt"
	"math"
	"math/cmplx"
)

func p(v []complex128) (r string) {
	for _, z := range v {
		r += pz(z)
	}
	return r
}
func pz(z complex128) string {
	r := cmplx.Abs(z)
	a := cmplx.Phase(z) * 180.0 / math.Pi
	if a < 0 {
		a += 360
	}
	return fmt.Sprintf("  %va%.5f", r, a)
}
func pm(m [][]complex128) {
	for _, v := range m {
		fmt.Println(p(v))
	}
}

type QR struct {
	H     [][]complex128
	Rdiag []complex128
	M, N  int
}

func NewQR(H [][]complex128) QR { // QR reusing H in column major format.
	m := len(H[0]) // Number of rows.
	n := len(H)    // Number of columns.
	if m < n {
		panic("qr: matrix is underdetermined")
	}
	Rdiag := make([]complex128, n)
	for j := 0; j < n; j++ {
		s := vectorNorm(H[j][j:])
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
		M:     m,
		N:     n,
	}
}
func (D QR) Solve(b []complex128) []complex128 {
	if len(b) != D.M {
		panic("qr: wrong input dimension for QR.Solve.")
	}
	return D.rSolve(D.qMul(b))
}
func (D QR) qMul(y []complex128) []complex128 {
	if len(y) != D.M {
		panic("qr: input vector lengths mismatch for qMul.")
	}
	//y := make([]complex128, D.M)
	//for i := 0; i < D.M; i++ {
	//	y[i] = x[i]
	//}
	for j := 0; j < D.N; j++ {
		var sum complex128
		for k := j; k < D.M; k++ {
			sum += cmplx.Conj(D.H[j][k]) * y[k]
		}
		for k := j; k < D.M; k++ {
			y[k] -= D.H[j][k] * sum
		}
	}
	fmt.Println("qmul y", p(y))
	return y
}
func (D QR) rSolve(b []complex128) []complex128 {
	if len(b) != D.M {
		panic("qr: input vector lengths mismatch for rSolve.")
	}
	x := make([]complex128, D.M)
	for i := 0; i < D.M; i++ {
		x[i] = b[i]
	}
	for i := D.N - 1; i >= 0; i-- {
		var sum complex128
		for j := i + 1; j < D.N; j++ {
			//x[i] -= D.H[j][i] * x[j]
			sum += D.H[j][i] * x[j]
		}
		fmt.Println("sum", sum)
		x[i] -= sum
		x[i] /= D.Rdiag[i]
	}
	return x[0:D.N]
}
func vectorNorm(x []complex128) (norm float64) {
	for i := 0; i < len(x); i++ {
		norm = math.Hypot(norm, cmplx.Abs(x[i]))
	}
	return
}

type RQ struct {
	H     [][]float64
	Rdiag []float64
	M, N  int
}

func NewRQ(H [][]float64) RQ { // column-major, overwrite
	m := len(H[0]) // Number of rows.
	n := len(H)    // Number of columns.
	if m < n {
		panic("qr: matrix is underdetermined")
	}
	Rdiag := make([]float64, n)
	for j := 0; j < n; j++ {
		s := RealNorm(H[j][j:])
		if s == 0 {
			panic("matrix contains zero-columns")
		}
		if H[j][j] > 0 {
			Rdiag[j] = -s
		} else {
			Rdiag[j] = s
		}
		f := 1.0 / math.Sqrt(s*(s+math.Abs(H[j][j])))
		H[j][j] -= Rdiag[j]
		for k := j; k < m; k++ {
			H[j][k] *= f
		}
		for i := j + 1; i < n; i++ {
			var sum float64
			for k := j; k < m; k++ {
				sum += H[j][k] * H[i][k]
			}
			for k := j; k < m; k++ {
				H[i][k] -= H[j][k] * sum
			}
		}
	}
	return RQ{
		H:     H,
		Rdiag: Rdiag,
		M:     m,
		N:     n,
	}
}
func (D RQ) Solve(b []float64) []float64 {
	if len(b) != D.M {
		panic("qr: wrong input dimension for QR.Solve.")
	}
	return D.rSolve(D.qMul(b))
}
func (D RQ) qMul(y []float64) []float64 {
	if len(y) != D.M {
		panic("qr: input vector lengths mismatch for qMul.")
	}
	//y := make([]float64, D.M)
	//for i := 0; i < D.M; i++ {
	//	y[i] = x[i]
	//}
	for j := 0; j < D.N; j++ {
		var sum float64
		for k := j; k < D.M; k++ {
			sum += D.H[j][k] * y[k]
		}
		for k := j; k < D.M; k++ {
			y[k] -= D.H[j][k] * sum
		}
	}
	return y
}
func (D RQ) rSolve(b []float64) []float64 {
	if len(b) != D.M {
		panic("qr: input vector lengths mismatch for rSolve.")
	}
	x := make([]float64, D.M)
	for i := 0; i < D.M; i++ {
		x[i] = b[i]
	}
	for i := D.N - 1; i >= 0; i-- {
		s := 0.0
		for j := i + 1; j < D.N; j++ {
			s += D.H[j][i] * x[j]
			x[i] -= D.H[j][i] * x[j]
		}
		x[i] /= D.Rdiag[i]
	}
	return x[0:D.N]
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
