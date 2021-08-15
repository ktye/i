package main

import (
	. "github.com/ktye/wg/module"
)

type rdf = func(K, int32, T, int32) K

func rd0(x K, yp int32, t T, n int32) K { return 0 }
func min(x K, yp int32, t T, n int32) K { // &/x
	xp := int32(x)
	switch t - 17 {
	case 0: // Bt
		if x == 0 {
			xp = 1
		}
		return Kb(mini(xp, all(yp, n)))
	case 1: // Ct
		if x == 0 {
			xp = 127
		}
		for i := int32(0); i < n; i++ {
			xp = mini(xp, I8(yp+i))
		}
		return Kc(xp)
	case 2: // It
		if x == 0 {
			xp = 2147483647
		}
		for i := int32(0); i < n; i++ {
			xp = mini(xp, I32(yp))
			yp += 4
		}
		return Ki(xp)
	case 3: // St
		if x == 0 {
			xp = (nn(K(I64(8))) << 3) - 8
		}
		for i := int32(0); i < n; i++ {
			xp = mini(xp, I32(yp))
			yp += 4
		}
		return Ks(xp)
	case 4: // Ft
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
	switch t - 17 {
	case 0: // Bt
		return Kb(maxi(xp, any(yp, n)))
	case 1: // Ct
		if x == 0 {
			xp = -128
		}
		for i := int32(0); i < n; i++ {
			xp = maxi(xp, I8(yp+i))
		}
		return Kc(xp)
	case 2: // It
		if x == 0 {
			xp = nai
		}
		for i := int32(0); i < n; i++ {
			xp = maxi(xp, I32(yp))
			yp += 4
		}
		return Ki(xp)
	case 3: // St
		for i := int32(0); i < n; i++ {
			xp = maxi(xp, I32(yp))
			yp += 4
		}
		return Ks(xp)
	case 4: // Ft
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
	switch t - 17 {
	case 0: // Bt
		return Ki(xp + sumb(yp, n))
	case 1: // Ct
		for i := int32(0); i < n; i++ {
			xp += I8(yp + i)
		}
		return Kc(xp)
	case 2: // It
		return Ki(xp + sumi(yp, n))
	case 3: // St
		return 0
	case 4: // Ft
		f := F64(xp)
		if x == 0 {
			f = 0.0
		}
		return Kf(f + sumf(yp, n, 8))
	case 5: // Zt
		var re, im float64
		if x != 0 {
			re, im = F64(xp), F64(xp+8)
		}
		return Kz(re+sumf(yp, n, 16), im+sumf(yp+8, n, 16))
	default:
		return 0
	}
}
func sumb(xp, xn int32) (r int32) {
	e := xp + xn
	ve := e &^ 7
	var s int64
	for xp < ve { // todo: I8x16popcnt when ready: https://github.com/WebAssembly/simd/pull/379
		s += I64popcnt(uint64(I64(xp)))
		xp += 8
	}
	r = int32(s)
	for xp < e {
		r += I8(xp)
		xp++
	}
	return r
}
func sumi(xp, xn int32) (r int32) {
	e := xp + 4*xn
	for xp < e {
		r += I32(xp)
		xp += 4
	}
	return r
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
	switch t - 17 {
	case 0: // Bt
		if x == 0 {
			xp = 1
		}
		return Kb(mini(xp, all(yp, n)))
	case 1: // Ct
		if x == 0 {
			xp = 1
		}
		for i := int32(0); i < n; i++ {
			xp *= I8(yp + i)
		}
		return Kc(xp)
	case 2: // It
		if x == 0 {
			xp = 1
		}
		for i := int32(0); i < n; i++ {
			xp *= I32(yp)
			yp += 4
		}
		return Ki(xp)
	case 3: // St
		return 0
	case 4: // Ft
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
	if t != Bt && t != It {
		return 0
	}
	r = mk(It, n)
	rp := int32(r)
	s := int32(x)
	e := yp
	if t == Bt {
		e += n
		for yp < e {
			s += I8(yp)
			SetI32(rp, s)
			rp += 4
			yp++
			continue
		}
	} else {
		e += 4 * n
		for yp < e {
			s += I32(yp)
			SetI32(rp, s)
			rp += 4
			yp += 4
			continue
		}
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
func mins(x K, yp int32, t T, n int32) K {
	if t != Bt {
		return 0
	} else {
		return 0
	}

}
func maxs(x K, yp int32, t T, n int32) K {
	if t != Bt || x != 0 {
		return 0
	} else {
		i := int32(fwh(yp, n))
		r := mk(Bt, n)
		rp := int32(r)
		if i < 0 {
			i = n
		}
		Memoryfill(rp, 0, i)
		Memoryfill(rp+i, 1, n-i)
		return r
	}
}
