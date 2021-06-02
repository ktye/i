package main

import . "github.com/ktye/wg/module"

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
		switch xt - 17 {
		case 0: // Bt
			return mtC(xp, yp, e)
		case 1: // Ct
			return mtC(xp, yp, e)
		case 2: // It
			return mtI(xp, yp, e)
		case 3: // St
			return mtI(xp, yp, e)
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
func mtC(xp, yp, e int32) (r int32) {
	m := maxi(yp, e-16)
	for yp < m {
		if I8x16load(xp).Eq(I8x16load(yp)).All_true() == 0 {
			return 0
		}
		xp += 16
		yp += 16
		continue
	}
	for yp < e {
		if I8(xp) != I8(yp) {
			return 0
		}
		xp++
		yp++
		continue
	}
	return 1
}
func mtI(xp, yp, e int32) (r int32) {
	m := maxi(yp, e-16)
	for yp < m {
		if I8x16load(xp).Eq(I8x16load(yp)).All_true() == 0 {
			return 0
		}
		xp += 16
		yp += 16
		continue
	}
	for yp < e {
		if I32(xp) != I32(yp) {
			return 0
		}
		xp += 4
		yp += 4
		continue
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
