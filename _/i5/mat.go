package k

import "github.com/ktye/k/mat"

func qr(x T) T {
	defer dx(x)
	a := matrix(x.(L))
	switch v := a.(type) {
	case [][]float64:
		return mat.NewRQ(v)
	case [][]complex128:
		return mat.NewQR(v)
	default:
		panic("type")
	}
}
func solve(x, y T) T {
	if _, o := y.(I); o {
		y = uptype(y, floattype)
	}
	switch v := x.(type) {
	case mat.RQ:
		if _, o := y.(L); o {
			return eachright(f2(solve), x, y)
		} else {
			y = use(y.(vector))
			return KF(v.Solve(y.(F).v))
		}
	case mat.QR:
		if _, o := y.(L); o {
			return eachright(f2(solve), x, y)
		} else {
			y = use(y.(vector))
			return KZ(v.Solve(y.(Z).v))
		}
	case L:
		return solve(qr(x), y)
	default:
		panic("type")
	}
}
func matrix(x L) T { // column major [][]float64, [][]complex128 or panic.
	if len(x.v) == 0 {
		dx(x)
		return make([][]float64, 0)
	}
	switch v := x.v[0].(type) {
	case I:
		return matrix(add(x, 0.0).(L))
	case F:
		return realmatrix(x, len(v.v))
	case Z:
		return complexmatrix(x, len(v.v))
	default:
		panic("type")
	}
}
func realmatrix(x L, n int) (r [][]float64) {
	r = make([][]float64, len(x.v))
	for i := range x.v {
		xi := x.v[i].(F)
		if len(xi.v) != n {
			panic("uniform")
		}
		r[i] = make([]float64, n)
		copy(r[i], xi.v)
	}
	dx(x)
	return r
}
func complexmatrix(x L, n int) (r [][]complex128) {
	r = make([][]complex128, len(x.v))
	for i := range x.v {
		xi := x.v[i].(Z)
		if len(xi.v) != n {
			panic("uniform")
		}
		r[i] = make([]complex128, n)
		copy(r[i], xi.v)
	}
	dx(x)
	return r
}
