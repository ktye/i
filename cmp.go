package main

import (
	. "github.com/ktye/wg/module"
)

func Mtc(x, y K) (r K) {
	r = Kb(match(x, y))
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
		ve := e &^ 15
		switch xt - 17 {
		case 0: // Bt
			return mtC(xp, yp, ve, e)
		case 1: // Ct
			return mtC(xp, yp, ve, e)
		case 2: // It
			return mtI(xp, yp, ve, e)
		case 3: // St
			return mtI(xp, yp, ve, e)
		case 4: // Ft
			return mtF(xp, yp, e)
		case 5: // Zt
			return mtF(xp, yp, e)
		case 6: // Lt
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
		return ib(xp == yp)
	}
	switch xt - ft {
	case 0: // ft
		return eqf(F64(xp), F64(yp))
	case 1: // zt
		return eqz(F64(xp), F64(xp+8), F64(yp), F64(yp+8))
	}
	return 0 // no deep comparision for comp, derived, proj, lambda
}
func mtC(xp, yp, ve, e int32) (r int32) {
	for yp < ve {
		if I8x16load(xp).Eq(I8x16load(yp)).All_true() == 0 {
			return 0
		}
		xp += 16
		yp += 16
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
func mtI(xp, yp, ve, e int32) (r int32) {
	for yp < ve {
		if I8x16load(xp).Eq(I8x16load(yp)).All_true() == 0 {
			return 0
		}
		xp += 16
		yp += 16
	}
	for yp < e {
		if I32(xp) != I32(yp) {
			return 0
		}
		xp += 4
		yp += 4
	}
	return 1
}
func mtF(xp, yp, e int32) (r int32) {
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
func all(x, n int32) int32 {
	e := x + n
	ve := e &^ 15
	for x < ve {
		if I8x16load(x).All_true() == 0 {
			return 0
		}
		x += 16
	}
	for x < e {
		if I8(x) != 0 {
			return 0
		}
		x++
	}
	return 1
}
func any(x, n int32) int32 {
	e := x + n
	ve := e &^ 15
	for x < ve {
		if I8x16load(x).Any_true() != 0 {
			return 1
		}
		x += 16
	}
	for x < e {
		if I8(x) != 0 {
			return 1
		}
		x++
	}
	return 0
}
func Any(x K) (r K) {
	if tp(x) != Bt {
		trap(Type)
	}
	xn := nn(x)
	dx(x)
	return Kb(any(int32(x), xn))
}
func In(x, y K) (r K) {
	xt, yt := tp(x), tp(y)
	if xt == yt && xt > 16 {
		return Ecl(30, l2(x, y))
	} else if xt+16 != yt {
		trap(Type)
	}
	r = in(x, y, xt)
	dx(y)
	return r
}
func in(x, y K, xt T) K {
	xp, yp := int32(x), int32(y)
	e := ep(y)
	ve := e &^ 15
	switch xt - 1 {
	case 0: //bt
		e = inC(xp, yp, ve, e)
	case 1: //ct
		e = inC(xp, yp, ve, e)
	case 2: //it
		e = inI(xp, yp, ve, e)
	case 3: //st
		e = inI(xp, yp, ve, e)
	case 4: //ft
		dx(x)
		e = inF(F64(xp), yp, e)
	case 5: //zt
		dx(x)
		e = inZ(F64(xp), F64(xp+8), yp, e)
	default:
		trap(Type)
	}
	return Kb(ib(e != 0))
}
func inC(x, yp, ve, e int32) int32 {
	v := I8x16splat(x)
	for yp < ve {
		if v.Eq(I8x16load(yp)).Any_true() != 0 {
			return yp
		}
		yp += 16
	}
	for yp < e {
		if x == I8(yp) {
			return yp
		}
		yp++
	}
	return 0
}
func inI(x, yp, ve, e int32) int32 {
	v := I32x4splat(x)
	for yp < ve {
		if v.Eq(I32x4load(yp)).Any_true() != 0 {
			return yp
		}
		yp += 16
	}
	for yp < e {
		if x == I32(yp) {
			return yp
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
	xt := tp(x)
	xp := int32(x)
	if xt == bt {
		r = Kb(1 - xp)
	} else if xt == Bt {
		r = use1(x)
		rp := int32(r)
		e := ep(r)
		w := I8x16splat(1)
		for rp < e {
			I8x16store(rp, I8x16load(xp).Not().And(w))
			xp += 16
			rp += 16
			continue
		}
		dx(x)
	} else if xt&15 == st {
		r = Eql(Ks(0), x)
	} else {
		r = Eql(Ki(0), x)
	}
	return r
}
