package main

// why? be kos
//
// f64:  +    -    *    /    =    !=   >    >=   <    <=
//      fadd fsub fmul fdiv feql fneq fgth fgte flth flte
//           fneg
// cnv: fi64 if64
//
// rewrite:
//  f64.load
//  f64.store
//  f64.convert_i32_s
//  f64.convert_i64_u
//  i64.trunc_f64_u
//  i32.trunc_f64_s
//
// nop:
//  f64.reinterpret_i64
//  i64.reinterpret_f64
//
// todo
//  f64.const
//  f64.abs
//  f64.max
//  f64.min
//  f64.floor
//  f64.copysign              go.dev/src/math/copysign.go
//  f64.sqrt
//
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
				m <<= uint64(1)
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
		m <<= uint64(1)
		e--
	}
	for m >= 18014398509481984 {
		t |= m & 1
		m >>= uint64(1)
		e++
	}
	if m >= 9007199254740992 {
		if m&1 != 0 && (t != 0 || m&2 != 0) {
			m++
			if m >= 18014398509481984 {
				m >>= uint64(1)
				e++
			}
		}
		m >>= uint64(1)
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
			m >>= uint64(1)
			e++
		}
		if m&1 != 0 && (t != 0 || m&2 != 0) {
			m++
		}
		m >>= uint64(1)
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

	if (fn | gn) != 0 {
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
	s := uint64(fe - ge)
	fm <<= uint64(2)
	gm <<= uint64(2)
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
func fsub(x, y uint64) uint64 { return fadd(x, fneg(y)) }
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
	if (fn | gn) != 0 {
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
func fdiv(x, y uint64) uint64 {
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
	if (fn | gn) != 0 {
		return nau
	}
	if (fi & gi) != 0 {
		return nau
	}
	if fi == 0 && gi == 0 && fm == 0 && gm == 0 {
		return nau
	}
	if fi != 0 || (gi == 0 && gm == 0) {
		return fs ^ gs ^ uint64(0x7FF0000000000000)
	}
	if gi != 0 || fm == 0 {
		return fs ^ gs ^ 0
	}

	//divlu 128/64->64quot,64rem
	var q, r uint64
	u1 := fm >> 10
	u0 := fm << 54
	v := gm
	b := uint64(4294967296)
	if u1 >= v {
		q = uint64(18446744073709551615)
		r = q
	} else {
		s := uint64(0)
		if v&uint64(9223372036854775808) == 0 {
			s++
			v <<= uint64(1)
		}
		vn1 := v >> 32
		vn0 := v & uint64(4294967295)
		un32 := u1<<s | u0>>(64-s)
		un10 := u0 << s
		un1 := un10 >> 32
		un0 := un10 & uint64(4294967295)
		q1 := un32 / vn1
		rh := un32 - q1*vn1

		l1 := int32(1)
		for l1 != 0 {
			if q1 >= b || q1*vn0 > b*rh+un1 {
				q1--
				rh += vn1
				l1 = I32B(rh < b)
			} else {
				l1 = 0
			}
		}

		un21 := un32*b + un1 - q1*v
		q0 := un21 / vn1
		rh = un21 - q0*vn1

		l1 = 1
		for l1 != 0 {
			if q0 >= b || q0*vn0 > b*rh+un0 {
				q0--
				rh += vn1
				l1 = I32B(rh < b)
			} else {
				l1 = 0
			}
		}
		q = q1*b + q0
		r = (un21*b + un0 - q0*v) >> s
	}
	return fpak(fs^gs, q, fe-ge-2, r)
}
func feql(x, y uint64) int32 {
	if (fnan(x) | fnan(y)) != 0 {
		return 0
	}
	return I32B(fcmp(x, y) == 0)
}
func fneq(x, y uint64) int32 {
	if (fnan(x) | fnan(y)) != 0 {
		return 1
	}
	return I32B(fcmp(x, y) != 0)
}
func fgth(x, y uint64) int32 {
	if (fnan(x) | fnan(y)) != 0 {
		return 0
	}
	return I32B(fcmp(x, y) > 0)
}
func fgte(x, y uint64) int32 {
	if (fnan(x) | fnan(y)) != 0 {
		return 0
	}
	return I32B(fcmp(x, y) >= 0)
}
func flte(x, y uint64) int32 {
	if (fnan(x) | fnan(y)) != 0 {
		return 0
	}
	return I32B(fcmp(x, y) <= 0)
}
func flth(x, y uint64) int32 {
	if (fnan(x) | fnan(y)) != 0 {
		return 0
	}
	return I32B(fcmp(x, y) < 0)
}
func fcmp(x, y uint64) int32 {
	fs := fsgn(x)
	fm := fmnt(x)
	fi := finf(x)
	gs := fsgn(y)
	gm := fmnt(y)
	gi := finf(y)
	if fi == 0 && gi == 0 && fm == 0 && gm == 0 {
		return 0
	}
	if fs > gs {
		return -1
	}
	if fs < gs {
		return 1
	}
	if (fs == 0 && x < y) || (fs != 0 && x > y) {
		return -1
	}
	if (fs == 0 && x > y) || (fs != 0 && x < y) {
		return 1
	}
	return 0
}
func fi64(x int64) uint64 { // f from i64
	f := uint64(x) & uint64(9223372036854775808)
	m := uint64(x)
	if f != 0 {
		m = -m
	}
	return fpak(f, m, 52, 0)
}
func if64(x uint64) int64 {
	fs := fsgn(x)
	fm := fmnt(x)
	fe := fexp(x)
	fi := finf(x)
	fn := fnan(x)
	if (fi | fn) != 0 {
		return 0
	}
	if fe < -1 {
		return 0
	}
	if fe > 63 {
		if fs != 0 && fm == 0 {
			return int64(-9223372036854775808)
		}
		if fs != 0 {
			return 0
		}
		return 0
	}
	for fe > 52 {
		fe--
		fm <<= uint64(1)
	}
	for fe < 52 {
		fe++
		fm >>= uint64(1)
	}
	r := int64(fm)
	if fs != 0 {
		r = -r
	}
	return r
}
func fabs(x uint64) uint64 { return x &^ 9223372036854775808 }
func flor(x uint64) uint64 {
	fs := fsgn(x)
	fm := fmnt(x)
	fi := finf(x)
	fn := fnan(x)
	if fm == 0 || (fn|fi) != 0 {
		return x
	}
	x &= uint64(9223372036854775807) //clear sign

	i := if64(x)
	r := fi64(i)
	if fs != 0 && x != r {
		r = fi64(i+1)
	}
	return r|fs

	/*
	x = fmod(x)
	if fs != 0 {
		// todo if not integral
		// x += 1.0
	}
	return x|fs
	*/
}
/*
func flor(x uint64) uint64 {
	fs := fsgn(x)
	fm := fmnt(x)
	fi := finf(x)
	fn := fnan(x)
	if fm == 0 || (fn|fi) != 0 {
		return x
	}
	x &= uint64(9223372036854775807) //clear sign
	x = fmod(x)
	if fs != 0 {
		// todo if not integral
		// x += 1.0
	}
	return x|fs
}
func fmod(x uint64) uint64 { // x>0.
	if x == 0 {
		return x
	}
	if x < uint64(4607182418800017408) { // x<1.0
		return 0
	}
	e := (x>>uint64(52))&2047 - 1023
	if e < 52 {
		x &^= uint64(1)<<(52-e) - 1
	}
	return x
}

func fmod(x uint64) uint64 {
	shift := 52
	mask  := 2047
	bias  := 1023
	if f < 1 {
		switch {
		case f < 0:
			int, frac = Modf(-f)
			return -int, -frac
		case f == 0:
			return f, f // Return -0, -0 when f == -0
		}
		return 0, f
	}

	x := Float64bits(f)
	e := uint(x>>shift)&mask - bias

	// Keep the top 12+e bits, the integer part; clear the rest.
	if e < 64-12 {
		x &^= 1<<(64-12-e) - 1
	}
	int = Float64frombits(x)
	frac = f - int
	return
}
*/

