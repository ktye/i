package main

import (
	. "github.com/ktye/wg/module"
)

// softfloat implementation of cosin_ atan2 log exp pow frexp is 2464 b

const pi float64 = 3.141592653589793
const maxfloat float64 = 1.797693134862315708145274237317043567981e+308

func hypot(p, q float64) float64 {
	p, q = F64abs(p), F64abs(q)
	if p < q {
		t := p
		p = q
		q = t
	}
	if p == 0.0 {
		return 0.0
	}
	q = q / p
	return p * F64sqrt(1+q*q)
}
func cosin(deg float64, rp int32) {
	var c, s float64
	if deg == 0 {
		c = 1.0
	} else if deg == 90 {
		s = 1.0
	} else if deg == 180 {
		c = -1.0
	} else if deg == 270 {
		s = -1.0
	} else {
		cosin_(deg*0.017453292519943295, rp, 0)
		return
	}
	SetF64(rp, c)
	SetF64(rp+8, s)
}
func ang2(y, x float64) (deg float64) {
	if y == 0 {
		if x < 0 {
			return 180.0
		}
		return 0.
	}
	if x == 0 {
		if y < 0 {
			return 270.0
		}
		return 90.0
	}
	deg = 57.29577951308232 * atan2(y, x)
	if deg < 0 {
		deg += 360.0
	}
	return deg
}
func cosin_(x float64, rp int32, csonly int32) {
	var c, s float64
	var ss, cs int32
	if x < 0 {
		x = -x
		ss = 1
	}
	var j uint64
	var y, z float64
	j = uint64(x * 1.2732395447351628) // *4/pi
	y = float64(j)
	if j&1 == 1 {
		j++
		y++
	}
	j &= 7
	z = ((x - y*7.85398125648498535156e-1) - y*3.77489470793079817668e-8) - y*2.69515142907905952645e-15
	if j > 3 {
		j -= 4
		//ss, cs = !ss, !cs
		ss, cs = 1-ss, 1-cs
	}
	if j > 1 {
		cs = 1 - cs
	}
	zz := z * z
	c = 1.0 - 0.5*zz + zz*zz*((((((-1.13585365213876817300e-11*zz)+2.08757008419747316778e-9)*zz+-2.75573141792967388112e-7)*zz+2.48015872888517045348e-5)*zz+-1.38888888888730564116e-3)*zz+4.16666666666665929218e-2)
	s = z + z*zz*((((((1.58962301576546568060e-10*zz)+-2.50507477628578072866e-8)*zz+2.75573136213857245213e-6)*zz+-1.98412698295895385996e-4)*zz+8.33333333332211858878e-3)*zz+-1.66666666666666307295e-1)
	if j == 1 || j == 2 {
		x = c
		c = s
		s = x
	}
	if cs != 0 {
		c = -c
	}
	if ss != 0 {
		s = -s
	}
	SetF64(rp, c)
	if csonly == 0 {
		SetF64(rp+8, s)
	} else if csonly == 1 {
		SetF64(rp, s)
	}
}
func atan2(y, x float64) float64 {
	// todo nan/inf
	q := atan(y / x)
	if x < 0 {
		if q <= 0 {
			return q + pi
		}
		return q - pi
	}
	return q
}
func atan(x float64) float64 {
	//if x == 0 {
	//	return x
	//}
	if x > 0 {
		return satan(x)
	} else {
		return -satan(-x)
	}
}
func satan(x float64) float64 {
	if x <= 0.66 {
		return xatan(x)
	}
	if x > 2.41421356237309504880 {
		return 1.5707963267948966 - xatan(1.0/x) + 6.123233995736765886130e-17
	}
	return 0.7853981633974483 + xatan((x-1)/(x+1)) + 0.5*6.123233995736765886130e-17
}
func xatan(x float64) float64 {
	z := x * x
	z = z * ((((-8.750608600031904122785e-01*z+-1.615753718733365076637e+01)*z+-7.500855792314704667340e+01)*z+-1.228866684490136173410e+02)*z + -6.485021904942025371773e+01) / (((((z+2.485846490142306297962e+01)*z+1.650270098316988542046e+02)*z+4.328810604912902668951e+02)*z+4.853903996359136964868e+02)*z + 1.945506571482613964425e+02)
	z = x*z + x
	return z
}
func exp(x float64) float64 {
	if x != x {
		return x
	}
	if x > 7.09782712893383973096e+02 {
		return inf
	}
	if x < -7.45133219101941108420e+02 {
		return 0.0
	}
	if -3.725290298461914e-09 < x && x < 3.725290298461914e-09 {
		return 1.0 + x
	}
	var k int64
	if x < 0 {
		k = int64(1.44269504088896338700*x - 0.5)
	} else {
		k = int64(1.44269504088896338700*x + 0.5)
	}
	hi := x - float64(k)*6.93147180369123816490e-01
	lo := float64(k) * 1.90821492927058770002e-10
	return expmulti(hi, lo, k)
}
func expmulti(hi, lo float64, k int64) float64 {
	r := hi - lo
	t := r * r
	c := r - t*(1.66666666666666657415e-01+t*(-2.77777777770155933842e-03+t*(6.61375632143793436117e-05+t*(-1.65339022054652515390e-06+t*4.13813679705723846039e-08))))
	y := 1 - ((lo - (r*c)/(2-c)) - hi)
	return ldexp(y, k)
}

//func signbit(x float64) int32 { return int32(I64reinterpret_f64(x) >> 63) }
func ldexp(frac float64, exp int64) float64 {
	if frac == 0 || frac > maxfloat || frac < -maxfloat || (frac != frac) {
		return frac
	}
	nf := normalize(frac)
	if nf != frac {
		exp -= 52
		frac = nf
	}
	x := uint64(I64reinterpret_f64(frac))
	exp += int64(x>>52)&2047 - 1023
	if exp < int64(-1075) {
		return F64copysign(0, frac)
	}
	if exp > int64(1023) {
		if frac < 0 {
			return -inf
		}
		return inf
	}
	var m float64 = 1.0
	if exp < int64(-1022) {
		exp += 53
		m = 1.1102230246251565e-16
	}
	x &^= 9218868437227405312
	x |= uint64(exp+1023) << 52
	return m * F64reinterpret_i64(uint64(x))
}
func frexp(f float64) (float64, int64) {
	var exp int64
	if f == 0.0 {
		return f, 0
	}
	if f < -maxfloat || f > maxfloat || (f != f) {
		return f, 0
	}
	nf := normalize(f)
	if nf != f {
		exp = -52
		f = nf
	}
	x := I64reinterpret_f64(f)
	exp += int64((x>>52)&2047) - 1022
	x &^= 9218868437227405312
	x |= 4602678819172646912
	return F64reinterpret_i64(x), exp
}
func normalize(x float64) float64 {
	if F64abs(x) < 2.2250738585072014e-308 {
		return x * 4.503599627370496e+15
	}
	return x
}
func log(x float64) float64 {
	if (x != x) || x > maxfloat {
		return x
	}
	if x < 0 {
		return na
	}
	if x == 0 {
		return -inf
	}
	f1, ki := frexp(x)
	if f1 < 0.7071067811865476 {
		f1 *= 2
		ki--
	}
	f := f1 - 1
	k := float64(ki)
	s := f / (2 + f)
	s2 := s * s
	s4 := s2 * s2
	t1 := s2 * (6.666666666666735130e-01 + s4*(2.857142874366239149e-01+s4*(1.818357216161805012e-01+s4*1.479819860511658591e-01)))
	t2 := s4 * (3.999999999940941908e-01 + s4*(2.222219843214978396e-01+s4*1.531383769920937332e-01))
	R := t1 + t2
	hfsq := 0.5 * f * f
	return k*6.93147180369123816490e-01 - ((hfsq - (s*(hfsq+R) + k*1.90821492927058770002e-10)) - f)
}
func modabsfi(f float64) float64 {
	if f < 1.0 {
		// simplified for f > 0
		return 0
	}
	x := I64reinterpret_f64(f)
	e := (x>>52)&2047 - 1023
	if e < 52 {
		x &^= 1<<(52-e) - 1
	}
	return F64reinterpret_i64(x)
}
func pow(x, y float64) float64 {
	if y == 0.0 || x == 1.0 {
		return 1.0
	}
	if y == 1.0 {
		return x
	}
	if (x != x) || (y != y) || y > maxfloat || y < -maxfloat { // simplified
		return na
	}
	if x == 0 { // simplified
		if y < 0 {
			return inf
		} else {
			return 0.0
		}
	}
	if y == 0.5 {
		return F64sqrt(x)
	}
	if y == -0.5 {
		return 1.0 / F64sqrt(x)
	}

	yf := F64abs(y)
	yi := modabsfi(yf)
	yf -= yi
	if yf != 0.0 && x < 0.0 {
		return na
	}
	if yi >= 9.223372036854776e+18 {
		if x == -1.0 {
			return 1.0
		} else if (F64abs(x) < 1.0) == (y > 0.0) {
			return 0.0
		} else {
			return inf
		}
	}
	a1 := 1.0
	ae := int64(0)
	if yf != 0 {
		if yf > 0.5 {
			yf -= 1.0
			yi += 1.0
		}
		a1 = exp(yf * log(x))
	}
	x1, xe := frexp(x)
	for i := int64(yi); i != 0; i >>= int64(1) {
		if xe < int64(-4096) || 4096 < xe {
			ae += xe
			break
		}
		if i&1 == 1 {
			a1 *= x1
			ae += xe
		}
		x1 *= x1
		xe <<= int64(1)
		if x1 < 0.5 {
			x1 += x1
			xe--
		}
	}
	if y < 0.0 {
		a1 = 1.0 / a1
		ae = -ae
	}
	return ldexp(a1, ae)
}
func ipow(x K, y int32) (r K) {
	if tp(x) == It {
		n := nn(x)
		r = mk(It, n)
		rp := int32(r)
		xp := int32(x)
		e := rp + 4*n
		for rp < e {
			SetI32(rp, iipow(I32(xp), y))
			xp += 4
			rp += 4
		}
		dx(x)
		return r
	} else {
		return Ki(iipow(int32(x), y))
	}
}
func iipow(x, y int32) (r int32) {
	r = 1
	for {
		if y&1 == 1 {
			r *= x
		}
		y >>= 1
		if y == 0 {
			break
		}
		x *= x
	}
	return r
}
