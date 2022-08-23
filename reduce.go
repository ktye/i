package main

import (
	. "github.com/ktye/wg/module"
)

func rd0(x K, yp int32, t T, n int32) K { return 0 }
func min(x K, yp int32, t T, n int32) K { // &/x
	xp := int32(x)
	switch t - 18 {
	case 0: // Ct
		if x == 0 {
			xp = 127
		}
		for i := int32(0); i < n; i++ {
			xp = mini(xp, I8(yp+i))
		}
		return Kc(xp)
	case 1: // It
		if x == 0 {
			xp = 2147483647
		}
		for i := int32(0); i < n; i++ {
			xp = mini(xp, I32(yp))
			yp += 4
		}
		return Ki(xp)
	case 2: // St
		if x == 0 {
			xp = (nn(K(I64(8))) << 3) - 8
		}
		for i := int32(0); i < n; i++ {
			xp = mini(xp, I32(yp))
			yp += 4
		}
		return Ks(xp)
	case 3: // Ft
		f := F64(xp)
		if x == 0 {
			f = F64reinterpret_i64(uint64(0x7FF0000000000000))
		}
		for i := int32(0); i < n; i++ {
			f = F64min(f, F64(yp))
			yp += 8
		}
		return Kf(f)
	default:
		return 0
	}
}
func max(x K, yp int32, t T, n int32) K { // |/x
	xp := int32(x)
	switch t - 18 {
	case 0: // Ct
		if x == 0 {
			xp = -128
		}
		for i := int32(0); i < n; i++ {
			xp = maxi(xp, I8(yp+i))
		}
		return Kc(xp)
	case 1: // It
		if x == 0 {
			xp = nai
		}
		for i := int32(0); i < n; i++ {
			xp = maxi(xp, I32(yp))
			yp += 4
		}
		return Ki(xp)
	case 2: // St
		for i := int32(0); i < n; i++ {
			xp = maxi(xp, I32(yp))
			yp += 4
		}
		return Ks(xp)
	case 3: // Ft
		f := F64(xp)
		if x == 0 {
			f = F64reinterpret_i64(uint64(0xFFF0000000000000))
		}
		for i := int32(0); i < n; i++ {
			f = F64max(f, F64(yp))
			yp += 8
		}
		return Kf(f)
	default:
		return 0
	}
}
func sum(x K, yp int32, t T, n int32) K { // +/x
	xp := int32(x)
	switch t - 18 {
	case 0: // Ct
		for i := int32(0); i < n; i++ {
			xp += I8(yp + i)
		}
		return Kc(xp)
	case 1: // It
		return Ki(xp + sumi(yp, n))
	case 2: // St
		return 0
	case 3: // Ft
		f := F64(xp)
		if x == 0 {
			f = 0.0
		}
		return Kf(f + sumf(yp, n, 8))
	case 4: // Zt
		var re, im float64
		if x != 0 {
			re, im = F64(xp), F64(xp+8)
		}
		return Kz(re+sumf(yp, n, 16), im+sumf(yp+8, n, 16))
	default:
		return 0
	}
}
func sumf(xp, n, s int32) (r float64) {
	if n < 128 {
		for i := int32(0); i < n; i++ {
			r += F64(xp)
			xp += s
		}
		return r
	}
	m := n / 2
	return sumf(xp, m, s) + sumf(xp+s*m, n-m, s)
}
func prd(x K, yp int32, t T, n int32) K { // */x
	xp := int32(x)
	switch t - 18 {
	case 0: // Ct
		if x == 0 {
			xp = 1
		}
		for i := int32(0); i < n; i++ {
			xp *= I8(yp + i)
		}
		return Kc(xp)
	case 1: // It
		if x == 0 {
			xp = 1
		}
		for i := int32(0); i < n; i++ {
			xp *= I32(yp)
			yp += 4
		}
		return Ki(xp)
	case 2: // St
		return 0
	case 3: // Ft
		f := F64(xp)
		if x == 0 {
			f = 1.0
		}
		for i := int32(0); i < n; i++ {
			f *= F64(yp)
			yp += 8
		}
		return Kf(f)
	default:
		return 0
	}
}

func sums(x K, yp int32, t T, n int32) (r K) {
	if t != It {
		return 0
	}
	r = mk(It, n)
	rp := int32(r)
	s := int32(x)
	e := yp + 4*n
	for yp < e {
		s += I32(yp)
		SetI32(rp, s)
		rp += 4
		yp += 4
		continue
	}
	return r
}
func prds(x K, yp int32, t T, n int32) (r K) {
	if t != It {
		return 0
	}
	r = mk(It, n)
	rp := int32(r)
	s := int32(x)
	if x == 0 {
		s = 1
	}
	e := yp + 4*n
	for yp < e {
		s *= I32(yp)
		SetI32(rp, s)
		rp += 4
		yp += 4
		continue
	}
	return r
}
