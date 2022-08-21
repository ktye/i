package main

import (
	. "github.com/ktye/wg/module"
)

func Mtc(x, y K) (r K) {
	r = Ki(match(x, y))
	dx(x)
	dx(y)
	return r
}
func match(x, y K) int32 {
	if x == y {
		return 1
	}
	xt, yt := tp(x), tp(y)
	if xt != yt {
		return 0
	}
	var xn, yn int32
	if xt > 16 {
		xn, yn = nn(x), nn(y)
		if xn != yn {
			return 0
		}
		if xn == 0 {
			return 1
		}
		xp, yp := int32(x), int32(y)
		e := ep(y)
		switch xt - 18 {
		case 0: // Ct
			return mtC(xp, yp, e)
		case 1: // It
			return mtC(xp, yp, e) //mtI
		case 2: // St
			return mtC(xp, yp, e) //mtI
		case 3: // Ft
			return mtF(xp, yp, e)
		case 4: // Zt
			return mtF(xp, yp, e)
		case 5: // Lt
			for i := int32(0); i < xn; i++ {
				if match(K(I64(xp)), K(I64(yp))) == 0 {
					return 0
				}
				xp += 8
				yp += 8
			}
			return 1
		default: // Dt, Tt
			if match(K(I64(xp)), K(I64(yp))) != 0 {
				return match(K(I64(xp+8)), K(I64(yp+8)))
			}
			return 0
		}
	}
	xp, yp := int32(x), int32(y)
	if xt < ft {
		return I32B(xp == yp)
	}
	switch int32(xt-ft) - 3*I32B(xt > 9) {
	case 0: // ft
		return eqf(F64(xp), F64(yp))
	case 1: // zt
		return eqz(F64(xp), F64(xp+8), F64(yp), F64(yp+8))
	case 2: // composition
		yn = 8 * nn(y)
	case 3: // derived
		yn = 16
	case 4: // projection
		yn = 24
	case 5: // lambda
		return match(K(I64(xp+16)), K(I64(yp+16))) // compare strings
	default: // xf
		return I32B(I64(xp) == I64(yp))
	}
	for yn > 0 { // composition, derived, projection
		yn -= 8
		if match(K(I64(xp+yn)), K(I64(yp+yn))) == 0 {
			return 0
		}
	}
	return 1
}
func mtC(xp, yp, e int32) int32 {
	ve := e &^ 7
	for yp < ve {
		if I64(xp) != I64(yp) {
			return 0
		}
		xp += 8
		yp += 8
	}
	for yp < e {
		if I8(xp) != I8(yp) {
			return 0
		}
		xp++
		yp++
	}
	return 1
}
func mtF(xp, yp, e int32) int32 {
	for yp < e {
		if eqf(F64(xp), F64(yp)) == 0 {
			return 0
		}
		xp += 8
		yp += 8
		continue
	}
	return 1
}
func In(x, y K) K {
	xt, yt := tp(x), tp(y)
	if xt == yt && xt > 16 {
		return Ecl(30, l2(x, y))
	} else if xt+16 != yt {
		trap(Type)
	}
	dx(y)
	return in(x, y, xt)
}
func in(x, y K, xt T) K {
	xp, yp := int32(x), int32(y)
	e := ep(y)
	switch xt - 2 {
	case 0: //ct
		e = inC(xp, yp, e)
	case 1: //it
		e = inI(xp, yp, e)
	case 2: //st
		e = inI(xp, yp, e)
	case 3: //ft
		dx(x)
		e = inF(F64(xp), yp, e)
	default: //zt
		dx(x)
		e = inZ(F64(xp), F64(xp+8), yp, e)
	}
	return Ki(I32B(e != 0))
}
func inC(x, yp, e int32) int32 {
	// maybe splat x to int64
	for yp < e {
		if x == I8(yp) {
			return yp //used in idxc
		}
		yp++
	}
	return 0
}
func inI(x, yp, e int32) int32 {
	for yp < e {
		if x == I32(yp) {
			return yp //used in idxi
		}
		yp += 4
	}
	return 0
}
func inF(x float64, yp int32, e int32) int32 {
	for yp < e {
		if eqf(x, F64(yp)) != 0 {
			return yp
		}
		yp += 8
	}
	return 0
}
func inZ(re, im float64, yp int32, e int32) int32 {
	for yp < e {
		if eqz(re, im, F64(yp), F64(yp+8)) != 0 {
			return yp
		}
		yp += 16
	}
	return 0
}

func Not(x K) (r K) { // ~x
	if tp(x)&15 == st {
		r = Eql(Ks(0), x)
	} else {
		r = Eql(Ki(0), x)
	}
	return r
}
