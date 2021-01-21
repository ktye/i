package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

// see github.com/ktye/{qr,svd}

// least squares (real and complex) A:column-major
// direct: x:solve[A;b]
// reuse:  x:solve[qr A;B]

func qr1(x uint32) uint32 {
	if zmatrix(x) {
		return K(NewQR(complexmatrix(x)))
	} else {
		return K(NewRQ(realmatrix(x)))
	}
}
func solve2(x, y uint32) (r uint32) {
	if tp(x) == 1 && nn(x) == 1 && MC[8+x] == 'c' {
		dx(x)
		return cond(y)
	}
	if tp(x) != 7 {
		x = qr1(x)
	}
	keys := Sk(MI[2+x>>2])
	if len(keys) != 4 || keys[1] != "Rdiag" {
		panic("solve2: no qr input")
	}
	t := tp(MI[2+1+MI[3+x>>2]>>2])
	if t == 3 {
		var q RQ
		G(x, &q)
		solve := func(b uint32) uint32 {
			return kF(q.Solve(Fk(b)))
		}
		if tp(y) == 6 {
			n := nn(y)
			r = mk(6, n)
			for i := uint32(0); i < n; i++ {
				MI[2+i+r>>2] = solve(MI[2+i+y>>2])
			}
		} else {
			r = solve(y)
		}
	} else {
		var q QR
		G(x, &q)
		solve := func(b uint32) uint32 {
			return kZ(q.Solve(Zk(b)))
		}
		if tp(y) == 6 {
			n := nn(y)
			r = mk(6, n)
			for i := uint32(0); i < n; i++ {
				MI[2+i+r>>2] = solve(MI[2+i+y>>2])
			}
		} else {
			r = solve(y)
		}
	}
	return dxyr(x, y, r)
}

// diag 5     -> 5x5 unit matrix int
// diag 1 2 3 -> 3x3 matrix (i|f|z)
// diag A     -> diagnal vector (r|z)
func diag1(x uint32) (r uint32) {
	t, n, xp := v1(x)
	if tp(x) == 2 && n == 1 {
		if m := MI[2+x>>2]; m > 1 {
			return dxr(x, diag1(take(mki(1), m)))
		}
	}
	if t == 6 {
		if nn(MI[2+x>>2]) != n {
			panic("diag: input is not square")
		}
		if zmatrix(x) {
			r = mk(4, n)
			for i := uint32(0); i < n; i++ {
				p := MI[2+i+x>>2]
				MF[1+2*i+r>>3] = MF[1+2*i+p>>3]
				MF[2+2*i+r>>3] = MF[2+2*i+p>>3]
			}
		} else {
			r = mk(3, n)
			for i := uint32(0); i < n; i++ {
				p := MI[2+i+x>>2]
				MF[1+i+r>>3] = MF[1+i+p>>3]
			}
		}
		return dxr(x, r)
	}
	if t < 2 || t > 4 {
		panic("diag:type")
	}
	r = mk(6, n)
	w := uint32(C(t))
	for i := uint32(0); i < n; i++ {
		ri := cst(mki(t), seq(0, n, 0))
		mv(8+ri+w*i, xp, w)
		MI[2+i+r>>2] = ri
		xp += w
	}
	return dxr(x, r)
}
func mul2(x, y uint32) (r uint32) { // mat*mat  mat*vec  XT*Y column-major f|z
	if tp(y) != 6 {
		y = enl(y)
		defer fst(r)
	}
	rows := nn(x)
	cols := nn(y)
	if nn(MI[2+x>>2]) != nn(MI[2+y>>2]) {
		panic("mul: inner dimensions")
	}
	r = mk(6, cols)
	if zmatrix(x) {
		if zmatrix(y) == false {
			panic("mul: expected complex Y")
		}
		for j := uint32(0); j < cols; j++ {
			rj := mk(4, rows)
			for i := uint32(0); i < rows; i++ {
				z := zvecvec(MI[2+i+x>>2], MI[2+j+y>>2])
				MF[1+2*i+rj>>3] = real(z)
				MF[2+2*i+rj>>3] = imag(z)
			}
			MI[2+j+r>>2] = rj
		}
	} else {
		if zmatrix(y) == true {
			panic("mul: expected real Y")
		}
		for j := uint32(0); j < cols; j++ {
			rj := mk(3, rows)
			for i := uint32(0); i < rows; i++ {
				MF[1+i+rj>>3] = rvecvec(MI[2+i+x>>2], MI[2+j+y>>2])
			}
			MI[2+j+r>>2] = rj
		}
	}
	return dxyr(x, y, r)
}
func zvecvec(x, y uint32) (r complex128) {
	n := nn(x)
	for i := uint32(0); i < n; i++ {
		r += complex(MF[1+2*i+x>>3], MF[2+2*i+x>>3]) * complex(MF[1+2*i+y>>3], MF[2+2*i+y>>3])
	}
	return r
}
func rvecvec(x, y uint32) (r float64) {
	n := nn(x)
	for i := uint32(0); i < n; i++ {
		r += MF[1+i+x>>3] * MF[1+i+y>>3]
	}
	return r
}
func cond(x uint32) uint32 {
	if zmatrix(x) == false {
		panic("cond: expect complex matrix")
	}
	s := NewSvd(complexmatrix(x))
	return kf(s.Condition())
}
func zmatrix(x uint32) (iscomplex bool) {
	n := nn(x)
	if t := tp(x); t != 6 {
		panic(fmt.Sprintf("expecting matrix, not %d/%d", t, n))
	}
	if n == 0 {
		panic("empty input matrix")
	}
	t := tp(MI[2+x>>2])
	m := nn(MI[2+x>>2])
	for i := uint32(0); i < n; i++ {
		p := MI[2+i+x>>2]
		mi := nn(p)
		ti := tp(p)
		if t < 3 || t > 4 || ti != t || mi != m {
			panic("matrix is not uniform")
		}
	}
	return t == 4
}
func complexmatrix(x uint32) [][]complex128 { // column major
	n := nn(x)
	r := make([][]complex128, n)
	for i := uint32(0); i < n; i++ {
		r[i] = Zk(MI[2+i+x>>2])
	}
	dx(x)
	return r
}
func realmatrix(x uint32) [][]float64 { // column major
	n := nn(x)
	r := make([][]float64, n)
	for i := uint32(0); i < n; i++ {
		r[i] = Fk(MI[2+i+x>>2])
	}
	dx(x)
	return r
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
func (D QR) qMul(x []complex128) []complex128 {
	if len(x) != D.M {
		panic("qr: input vector lengths mismatch for qMul.")
	}
	y := make([]complex128, D.M)
	for i := 0; i < D.M; i++ {
		y[i] = x[i]
	}
	for j := 0; j < D.N; j++ {
		var sum complex128
		for k := j; k < D.M; k++ {
			sum += cmplx.Conj(D.H[j][k]) * y[k]
		}
		for k := j; k < D.M; k++ {
			y[k] -= D.H[j][k] * sum
		}
	}
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
		for j := i + 1; j < D.N; j++ {
			x[i] -= D.H[j][i] * x[j]
		}
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
func (D RQ) qMul(x []float64) []float64 {
	if len(x) != D.M {
		panic("qr: input vector lengths mismatch for qMul.")
	}
	y := make([]float64, D.M)
	for i := 0; i < D.M; i++ {
		y[i] = x[i]
	}
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

// Svd is the result of a singular value decomposition A = U*diag(S)*conj(V).
// For a given input matrix A of size m x n, only a compact form is stored.
type Svd struct {
	U [][]complex128 // m x n
	S []float64      // Singular values, length n
	V [][]complex128 // n x n
}

// Condition returns the condition number of the original matrix.
func (s *Svd) Condition() float64 {
	if len(s.S) < 1 {
		return 0
	}
	return s.S[0] / s.S[len(s.S)-1]
}

// NewSvd computes the singular value decomposition of A.
// A = U*S*V'.
func NewSvd(H [][]complex128) Svd { // H: column major
	const eta = 2.8e-16
	const tol = 4.0e-293
	const zero complex128 = complex(0, 0)
	const one complex128 = complex(1, 0)

	norm := func(z complex128) float64 { return real(z)*real(z) + imag(z)*imag(z) }
	var b, c, t []float64
	var sn, w, x, y, z, cs, eps, f, g, h float64
	var i, j, k, k1, L, L1 int
	var q complex128
	var U, V [][]complex128
	var S []float64
	n := len(H)
	if n < 1 {
		panic("svd: input has no columns")
	}
	m := len(H[0])
	if m < 1 {
		panic("svd: matrix a has no rows")
	}
	for _, v := range H {
		if len(v) != m {
			panic("svd: input is not a uniform matrix")
		}
	}
	if m < n {
		panic("svd: input matrix has less rows than cols")
	}
	b = make([]float64, n)
	c = make([]float64, n)
	t = make([]float64, n)
	U = make([][]complex128, m)
	for i = range U {
		U[i] = make([]complex128, n)
	}
	S = make([]float64, n)
	V = make([][]complex128, n)
	for i = range V {
		V[i] = make([]complex128, n)
	}
	for {
		k1 = k + 1
		z = 0.0
		for i = k; i < m; i++ {
			z += norm(H[k][i])
		}
		b[k] = 0.0
		if z > tol {
			z = math.Sqrt(z)
			b[k] = z
			w = cmplx.Abs(H[k][k])
			q = one
			if w != 0.0 {
				q = H[k][k] / complex(w, 0)
			}
			H[k][k] = q * complex(z+w, 0)
			if k != n-1 {
				for j = k1; j < n; j++ {
					q = zero
					for i = k; i < m; i++ {
						q += cmplx.Conj(H[k][i]) * H[j][i]
					}
					q /= complex(z*(z+w), 0)
					for i = k; i < m; i++ {
						H[j][i] -= q * H[k][i]
					}
				}
			}
			q = -cmplx.Conj(H[k][k]) / complex(cmplx.Abs(H[k][k]), 0)
			for j = k1; j < n; j++ {
				H[j][k] *= q
			}
		}
		if k == n-1 {
			break
		}
		z = 0.0
		for j = k1; j < n; j++ {
			z += norm(H[j][k])
		}
		c[k1] = 0.0
		if z > tol {
			z = math.Sqrt(z)
			c[k1] = z
			w = cmplx.Abs(H[k1][k])
			q = one
			if w != 0.0 {
				q = H[k1][k] / complex(w, 0)
			}
			H[k1][k] = q * complex(z+w, 0)
			for i = k1; i < m; i++ {
				q = zero
				for j = k1; j < n; j++ {
					q += cmplx.Conj(H[j][k]) * H[j][i]
				}
				q /= complex(z*(z+w), 0)
				for j = k1; j < n; j++ {
					H[j][i] -= q * H[j][k]
				}
			}
			q = -cmplx.Conj(H[k1][k]) / complex(cmplx.Abs(H[k1][k]), 0)
			for i = k1; i < m; i++ {
				H[k1][i] *= q
			}
		}
		k = k1
	}
	eps = 0.0
	for k = 0; k < n; k++ {
		S[k] = b[k]
		t[k] = c[k]
		if S[k]+t[k] > eps {
			eps = S[k] + t[k]
		}
	}
	eps *= eta
	for j = 0; j < n; j++ {
		U[j][j] = one
		V[j][j] = one
	}
	for k = n - 1; k >= 0; k-- {
		for {
			for L = k; L >= 0; L-- {
				if math.Abs(t[L]) <= eps {
					goto Test
				}
				if math.Abs(S[L-1]) <= eps {
					break
				}
			}
			cs = 0.0
			sn = 1.0
			L1 = L - 1
			for i = L; i <= k; i++ {
				f = sn * t[i]
				t[i] *= cs
				if math.Abs(f) <= eps {
					goto Test
				}
				h = S[i]
				w = math.Sqrt(f*f + h*h)
				S[i] = w
				cs = h / w
				sn = -f / w
				for j = 0; j < n; j++ {
					x = real(U[j][L1])
					y = real(U[j][i])
					U[j][L1] = complex(x*cs+y*sn, 0)
					U[j][i] = complex(y*cs-x*sn, 0)
				}
			}
		Test:
			w = S[k]
			if L == k {
				break
			}
			x = S[L]
			y = S[k-1]
			g = t[k-1]
			h = t[k]
			f = ((y-w)*(y+w) + (g-h)*(g+h)) / (2.0 * h * y)
			g = math.Sqrt(f*f + 1.0)
			if f < 0.0 {
				g = -g
			}
			f = ((x-w)*(x+w) + (y/(f+g)-h)*h) / x
			cs = 1.0
			sn = 1.0
			L1 = L + 1
			for i = L1; i <= k; i++ {
				g = t[i]
				y = S[i]
				h = sn * g
				g = cs * g
				w = math.Sqrt(h*h + f*f)
				t[i-1] = w
				cs = f / w
				sn = h / w
				f = x*cs + g*sn
				g = g*cs - x*sn
				h = y * sn
				y = y * cs
				for j = 0; j < n; j++ {
					x = real(V[j][i-1])
					w = real(V[j][i])
					V[j][i-1] = complex(x*cs+w*sn, 0)
					V[j][i] = complex(w*cs-x*sn, 0)
				}
				w = math.Sqrt(h*h + f*f)
				S[i-1] = w
				cs = f / w
				sn = h / w
				f = cs*g + sn*y
				x = cs*y - sn*g
				for j = 0; j < n; j++ {
					y = real(U[j][i-1])
					w = real(U[j][i])
					U[j][i-1] = complex(y*cs+w*sn, 0)
					U[j][i] = complex(w*cs-y*sn, 0)
				}
			}
			t[L] = 0.0
			t[k] = f
			S[k] = x
		}
		if w >= 0.0 {
			continue
		}
		S[k] = -w
		for j = 0; j < n; j++ {
			V[j][k] = -V[j][k]
		}
	}
	for k = 0; k < n; k++ {
		g = -1.0
		j = k
		for i = k; i < n; i++ {
			if S[i] <= g {
				continue
			}
			g = S[i]
			j = i
		}
		if j == k {
			continue
		}
		S[j] = S[k]
		S[k] = g
		for i = 0; i < n; i++ {
			q = V[i][j]
			V[i][j] = V[i][k]
			V[i][k] = q
		}
		for i = 0; i < n; i++ {
			q = U[i][j]
			U[i][j] = U[i][k]
			U[i][k] = q
		}
	}
	for k = n - 1; k >= 0; k-- {
		if b[k] == 0.0 {
			continue
		}
		q = -H[k][k] / complex(cmplx.Abs(H[k][k]), 0)
		for j = 0; j < n; j++ {
			U[k][j] *= q
		}
		for j = 0; j < n; j++ {
			q = zero
			for i = k; i < m; i++ {
				q += cmplx.Conj(H[k][i]) * U[i][j]
			}
			q /= complex(cmplx.Abs(H[k][k])*b[k], 0)
			for i = k; i < m; i++ {
				U[i][j] -= q * H[k][i]
			}
		}
	}
	if n > 1 {
		for k = n - 2; k >= 0; k-- {
			k1 = k + 1
			if c[k1] == 0.0 {
				continue
			}
			q = -cmplx.Conj(H[k1][k]) / complex(cmplx.Abs(H[k1][k]), 0)
			for j = 0; j < n; j++ {
				V[k1][j] *= q
			}
			for j = 0; j < n; j++ {
				q = zero
				for i = k1; i < n; i++ {
					q += H[i][k] * V[i][j]
				}
				q /= complex(cmplx.Abs(H[k1][k])*c[k1], 0)
				for i = k1; i < n; i++ {
					V[i][j] -= q * cmplx.Conj(H[i][k])
				}
			}
		}
	}
	return Svd{U: U, S: S, V: V}
}
