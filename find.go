package main

import (
	. "github.com/ktye/wg/module"
)

type f3i = func(int32, int32, int32) int32

func Fnd(x, y K) K { // x?y
	r := K(0)
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		if yt == Tt {
			return grp(x, y)
		} else {
			return deal(x, y)
		}
	}
	if xt > Lt {
		if xt == Tt {
			trap() //nyi t?..
		}
		r = x0(x)
		return Atx(r, Fnd(r1(x), y))
	} else if xt == yt {
		return Ecr(18+16*K(I32B(yt == Lt)), l2(x, y))
	} else if xt == yt+16 {
		r = Ki(fnd(x, y, yt))
	} else if xt == Lt {
		return fdl(x, y)
	} else if yt == Lt {
		return Ecr(18, l2(x, y))
	} else {
		trap() //type
	}
	dxy(x, y)
	return r
}
func fnd(x, y K, t T) int32 {
	if nn(x) == 0 {
		return nai
	}
	xp := int32(x)
	r := Func[268+t].(f3i)(int32(y), xp, ep(x))
	if r == 0 {
		return nai
	}
	return (r - xp) >> (31 - I32clz(sz(16+t)))
}
func fdl(x, y K) K {
	xp := int32(x)
	dxy(x, y)
	e := ep(x)
	for xp < e {
		if match(K(I64(xp)), y) != 0 {
			return Ki((xp - int32(x)) >> 3)
		}
		xp += 8
	}
	return Ki(nai)
}
func idx(x, a, b int32) int32 {
	for i := a; i < b; i++ {
		if x == I8(i) {
			return i - a
		}
	}
	return -1
}

func Find(x, y K) K { // find[pattern;string] returns all matches (It)
	t := tp(x)
	if t != tp(y) || t != Ct {
		trap() //type
	}
	xn, yn := nn(x), nn(y)
	if xn*yn == 0 {
		dxy(x, y)
		return mk(It, 0)
	}
	r := mk(It, 0)
	yp := int32(y)
	e := yp + yn + 1 - xn
	for yp < e { // todo rabin-karp / knuth-morris / boyes-moore..
		if findat(int32(x), yp, xn) != 0 {
			r = cat1(r, Ki(yp-int32(y)))
			yp += xn
		} else {
			yp++
		}
		continue
	}
	dxy(x, y)
	return r
}
func findat(xp, yp, n int32) int32 {
	for i := int32(0); i < n; i++ {
		if I8(xp+i) != I8(yp+i) {
			return 0
		}
	}
	return 1
}

func Mtc(x, y K) K {
	dxy(x, y)
	return Ki(match(x, y))
}
func match(x, y K) int32 {
	if x == y {
		return 1
	}
	xt := tp(x)
	if xt != tp(y) {
		return 0
	}
	if xt > 16 {
		n := nn(x)
		if n != nn(y) {
			return 0
		}
		if n == 0 {
			return 1
		}
		xp, yp := int32(x), int32(y)
		if xt < Dt {
			return Func[246+xt].(f3i)(xp, yp, ep(y))
		} else {
			if match(K(I64(xp)), K(I64(yp))) != 0 {
				return match(K(I64(xp+8)), K(I64(yp+8)))
			}
			return 0
		}
	}
	yn := int32(0)
	xp, yp := int32(x), int32(y)
	if xt < ft {
		return I32B(xp == yp)
	}
	switch int32(xt-ft) - 3*I32B(xt > 9) {
	case 0: // ft
		return I32B(0 == cmF(xp, yp))
	case 1: // zt
		return I32B(0 == cmZ(xp, yp))
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
		if cmF(xp, yp) != 0 {
			return 0
		}
		xp += 8
		yp += 8
		continue
	}
	return 1
}
func mtL(xp, yp, e int32) int32 {
	for yp < e {
		if match(K(I64(xp)), K(I64(yp))) == 0 {
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
		trap() //type
	}
	dxy(x, y)
	return Ki(I32B(Func[268+xt].(f3i)(int32(x), int32(y), ep(y)) != 0))
}
func inC(x, yp, e int32) int32 {
	for yp < e { // maybe splat x to int64
		if x == I8(yp) {
			return yp
		}
		yp++
	}
	return 0
}
func inI(x, yp, e int32) int32 {
	for yp < e {
		if x == I32(yp) {
			return yp
		}
		yp += 4
	}
	return 0
}
func inF(xp, yp, e int32) int32 {
	for yp < e {
		if cmF(xp, yp) == 0 {
			return yp
		}
		yp += 8
	}
	return 0
}
func inZ(xp, yp, e int32) int32 {
	for yp < e {
		if cmZ(xp, yp) == 0 {
			return yp
		}
		yp += 16
	}
	return 0
}
