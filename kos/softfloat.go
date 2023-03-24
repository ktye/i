package main

// why? be kos

// from go.dev/src/runtime/softfloat64.go

const nau = uint64(0x7FF8000000000001)

func fsgn(x uint64) uint64 { return x & 9223372036854775808 }
func fmnt(x uint64) uint64 { return fupk(x, 1) }
func fexp(x uint64) int32  { return int32(fupk(x, 0)) }
func fnan(x uint64) int32  { return fnai(x, 1) }
func finf(x uint64) int32  { return fnai(x, 0) }
func fnai(x uint64, n int32) int32 {
	m := x & 4503599627370495
	e := int32(x>>52) & 2047
	if e == 2047 {
		if m != 0 {
			return n * 1
		}
		return 1 - n
	}
	return 0
}
func fupk(x uint64, mnt int32) uint64 {
	m := x & 4503599627370495
	e := int32(x>>52) & 2047
	if e == 2047 {
		return 0
	}
	if e == 0 {
		if m != 0 {
			e += -1022
			for m < 4503599627370496 {
				m <<= 1
				e--
			}
		}
	} else {
		m |= 4503599627370496
		e += -1023
	}
	if mnt != 0 {
		return m
	}
	return uint64(e)
}

func fpak(s, m uint64, e int32, t uint64) uint64 {
	//fmt.Println(" fpak", s, m, e, t)
	m0 := m
	e0 := e
	t0 := t
	if m == 0 {
		return s
	}
	for m < 4503599627370496 {
		m <<= 1
		e--
	}
	for m >= 18014398509481984 {
		t |= m & 1
		m >>= 1
		e++
	}
	if m >= 9007199254740992 {
		if m&1 != 0 && (t != 0 || m&2 != 0) {
			m++
			if m >= 18014398509481984 {
				m >>= 1
				e++
			}
		}
		m >>= 1
		e++
	}
	if e >= 1024 {
		return s ^ 9218868437227405312
	}
	if e < -1022 {
		if e < -1075 {
			return s | 0
		}
		m = m0
		e = e0
		t = t0
		for e < -1023 {
			t |= m & 1
			m >>= 1
			e++
		}
		if m&1 != 0 && (t != 0 || m&2 != 0) {
			m++
		}
		m >>= 1
		e++
		if m < 4503599627370496 {
			return s | m
		}
	}
	return s | uint64(e+1023)<<52 | m&(4503599627370495)
}
func fneg(x uint64) uint64 { return x ^ 9223372036854775808 }
func fadd(x, y uint64) uint64 {
	var ti int32
	var tu uint64
	fs := fsgn(x)
	fm := fmnt(x)
	fe := fexp(x)
	fi := finf(x)
	fn := fnan(x)
	gs := fsgn(y)
	gm := fmnt(y)
	ge := fexp(y)
	gi := finf(y)
	gn := fnan(y)

	if fn != 0 || gn != 0 {
		return nau
	}
	if fi != 0 && gi != 0 && fs != gs {
		return nau
	}
	if fi != 0 {
		return x
	}
	if gi != 0 {
		return y
	}
	if fm == 0 && gm == 0 && fs != 0 && gs != 0 {
		return x
	}
	if fm == 0 {
		if gm == 0 {
			y ^= gs
		}
		return y
	}
	if gm == 0 {
		return x
	}
	if fe < ge || fe == ge && fm < gm {
		tu = x
		x = y
		y = tu
		tu = fm
		fm = gm
		gm = tu
		tu = fs
		fs = gs
		gs = tu
		ti = fe
		fe = ge
		ge = ti
	}
	s := uint32(fe - ge)
	fm <<= 2
	gm <<= 2
	t := gm & (1<<s - 1)
	gm >>= s
	if fs == gs {
		fm += gm
	} else {
		fm -= gm
		if t != 0 {
			fm--
		}
	}
	if fm == 0 {
		fs = 0
	}
	//fmt.Println(" fs/fm/fe-2/t", fs, fm, fe-2, t)
	return fpak(fs, fm, fe-2, t)
}
func fsub(x, y uint64) uint64 { return fadd64(x, fneg64(y)) }
func fmul(x, y uint64) uint64 {
	fs := fsgn(x)
	fm := fmnt(x)
	fe := fexp(x)
	fi := finf(x)
	fn := fnan(x)
	gs := fsgn(y)
	gm := fmnt(y)
	ge := fexp(y)
	gi := finf(y)
	gn := fnan(y)
	if fn != 0 || gn != 0 {
		return nau
	}
	if fi != 0 && gi != 0 {
		return x ^ gs
	}
	if (fi != 0 && gm == 0) || (fm == 0 && gi != 0) {
		return nau
	}
	if fm == 0 {
		return x ^ gs
	}
	if gm == 0 {
		return y ^ fs
	}

	// mullu (multiply 64x64->128)
	u0 := fm & 4294967295
	u1 := fm >> 32
	v0 := gm & 4294967295
	v1 := gm >> 32
	w0 := u0 * v0
	t := u1*v0 + w0>>32
	w1 := t & 4294967295
	w2 := t >> 32
	w1 += u0 * v1
	lo := fm * gm
	hi := u1*v1 + w2 + w1>>32
	//fmt.Println("fmul   hi/lo", hi, lo, "fm/gm", fm, gm)

	tr := lo & 2251799813685247
	m := hi<<13 | lo>>51
	//fmt.Println("fmul", fs, gs, m, fe, ge, tr)
	return fpak(fs^gs, m, fe+ge-1, tr)
}

/*
func fneg(x float64) float64 { return 0 }
func fsqt(x float64) float64 { return 0 }
func fadd(x, y float64) float64 { return 0 }
func fsub(x, y float64) float64 { return 0 }
func fmul(x, y float64) float64 { return 0 }
func fdiv(x, y float64) float64 { return 0 }


*/

/*
mantbits64 uint = 52
expbits64  uint = 11
bias64          = -1<<(expbits64-1) + 1
nan64 uint64 = (1<<expbits64-1)<<mantbits64 + 1<<(mantbits64-1) // quiet NaN, 0 payload
inf64 uint64 = (1<<expbits64 - 1) << mantbits64
neg64 uint64 = 1 << (expbits64 + mantbits64)
*/
