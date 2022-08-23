//go:build small

package main

import (
	. "github.com/ktye/wg/module"
)

func nm(f int32, x K) (r K) { //monadic
	xt := tp(x)
	if xt > 16 {
		ff := K(32) //abs
		if f == 300 {
			ff = 5 //% sqrt
		}
		if f == 220 {
			ff = 3 //-
		}
		return Ech(ff, l1(x))
	}
	xp := int32(x)
	switch xt - 2 {
	case 0:
		return Kc(Func[f].(f1i)(xp))
	case 1:
		return Ki(Func[f].(f1i)(xp))
	case 2:
		return trap(Type)
	case 3:
		r = Kf(Func[1+f].(f1f)(F64(xp)))
		dx(x)
		return r
	case 4:
		r = Kz(Func[2+f].(f1z)(F64(xp), F64(xp+8)))
		dx(x)
		return r
	default:
		return trap(Type)
	}
}
func nd(f, ff int32, x, y K) (r K) { //dyadic
	var av int32
	var t T
	kf := K(ff)
	r, t, x, y = dctypes(x, y)
	if r != 0 {
		return key(r, Func[64+ff].(f2)(x, y), t)
	}
	av, t = conform(x, y)
	if av == 0 && t != lt {
		x, y = uptypes(x, y)
		t = tp(x)
		xp, yp := int32(x), int32(y)
		if ff > 7 && ff < 11 { // compare
			switch t - 2 {
			case 0: // ct
				return Ki(Func[f].(f2i)(xp, yp))
			case 1: // it
				return Ki(Func[f].(f2i)(xp, yp))
			case 2: // st
				return Ki(Func[f].(f2i)(xp, yp))
			case 3:
				dx(x)
				dx(y)
				return Ki(Func[1+f].(f2c)(F64(xp), F64(yp)))
			default:
				dx(x)
				dx(y)
				return Ki(Func[2+f].(f2d)(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
			}
		}
		switch t - 2 {
		case 0: // ct
			return Kc(Func[f].(f2i)(xp, yp))
		case 1: // it
			return Ki(Func[f].(f2i)(xp, yp))
		case 2: // st
			return trap(Type)
		case 3:
			dx(x)
			dx(y)
			return Kf(Func[1+f].(f2f)(F64(xp), F64(yp)))
		default:
			dx(x)
			dx(y)
			return Kz(Func[2+f].(f2z)(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		}
	}

	xy := l2(x, y)
	if av == 1 { //atom-vector
		return Ecr(kf, xy)
	}
	if av == 2 { // vector-vector
		return Ech(kf, xy) //vector-vector
	}
	return Ecl(kf, xy) // vector-atom (only comparisons)
}
func nc(f, ff int32, x, y K) K { return nd(f, ff, x, y) }

func Abs(x K) (r K) {
	r = nm(227, x)
	if tp(r)&15 == zt {
		r = Flr(r)
	}
	return r
}
func absz(x, y float64) (float64, float64) { return hypot(x, y), 0.0 }
func Mul(x, y K) K                         { return nd(256, 4, x, y) }
func Div(x, y K) K                         { return nd(267, 5, x, y) }
func Mod(x, y K) K                         { return nd(270, 41, x, y) }
func modi(x, y int32) int32                { return x % y }
func idiv(x, y K, mod int32) K {
	if mod != 0 {
		return Mod(x, y)
	}
	return Div(x, y)
}

func Sub(x, y K) K                                   { return nd(245, 3, x, y) }
func subz(xr, xi, yr, yi float64) (float64, float64) { return xr - yr, xi - yi }
