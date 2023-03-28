package main

// why? be kos
//
// f64:  +    -    *    /    =    !=   >    >=   <    <=
//      fadd fsub fmul fdiv feql fneq fmor fgte fles flte
//           fneg
//
//      fabs flor fcps(copysign) fmax fmin fsqr
//	fjcst     jfcst
//
// from go.dev/src/runtime/softfloat64.go

const nau = uint64(0x7FF8000000000001)

func fsgn(x uint64) uint64 { return x & uint64(0x8000000000000000) }
func fmnt(x uint64) uint64 { return fupk(x, 1) }
func fexp(x uint64) int32  { return int32(fupk(x, 0)) }
func fnan(x uint64) int32  { return fnai(x, 1) }
func finf(x uint64) int32  { return fnai(x, 0) }
func fnai(x uint64, n int32) int32 {
	m := x & uint64(0xfffffffffffff)
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
	m := x & uint64(0xfffffffffffff)
	e := int32(x>>52) & 2047
	if e == 2047 {
		return 0
	}
	if e == 0 {
		if m != 0 {
			e += -1022
			for m < uint64(0x10000000000000) {
				m <<= uint64(1)
				e--
			}
		}
	} else {
		m |= uint64(0x10000000000000)
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
	for m < uint64(0x10000000000000) {
		m <<= uint64(1)
		e--
	}
	for m >= uint64(0x40000000000000) {
		t |= m & 1
		m >>= uint64(1)
		e++
	}
	if m >= uint64(0x20000000000000) {
		if m&1 != 0 && (t != 0 || m&2 != 0) {
			m++
			if m >= uint64(0x40000000000000) {
				m >>= uint64(1)
				e++
			}
		}
		m >>= uint64(1)
		e++
	}
	if e >= 1024 {
		return s ^ uint64(0x7ff0000000000000)
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
		if m < uint64(0x10000000000000) {
			return s | m
		}
	}
	return s | uint64(e+1023)<<52 | m&uint64(0xfffffffffffff)
}
func fneg(x uint64) uint64 { return x ^ uint64(0x8000000000000000) }
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
	u0 := fm & uint64(0xffffffff)
	u1 := fm >> 32
	v0 := gm & uint64(0xffffffff)
	v1 := gm >> 32
	w0 := u0 * v0
	t := u1*v0 + w0>>32
	w1 := t & uint64(0xffffffff)
	w2 := t >> 32
	w1 += u0 * v1
	lo := fm * gm
	hi := u1*v1 + w2 + w1>>32
	//fmt.Println("fmul   hi/lo", hi, lo, "fm/gm", fm, gm)

	tr := lo & uint64(0x7ffffffffffff)
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
	b := uint64(0x100000000)
	if u1 >= v {
		q = uint64(0xffffffffffffffff)
		r = q
	} else {
		s := uint64(0)
		if v&uint64(0x8000000000000000) == 0 {
			s++
			v <<= uint64(1)
		}
		vn1 := v >> 32
		vn0 := v & uint64(0xffffffff)
		un32 := u1<<s | u0>>(64-s)
		un10 := u0 << s
		un1 := un10 >> 32
		un0 := un10 & uint64(0xffffffff)
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
func fmor(x, y uint64) int32 {
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
func fles(x, y uint64) int32 {
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
func fjcst(x int64) uint64 { // f from i64
	f := uint64(x) & uint64(0x8000000000000000)
	m := uint64(x)
	if f != 0 {
		m = -m
	}
	return fpak(f, m, 52, 0)
}
func jfcst(x uint64) int64 {
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
func ifcst(x uint64) int32 { return int32(jfcst(x)) }
func ficst(x int32) uint64 { return fjcst(int64(x)) }
func fabs(x uint64) uint64 { return x &^ uint64(0x8000000000000000) }
func flor(x uint64) uint64 {
	fs := fsgn(x)
	fm := fmnt(x)
	fi := finf(x)
	fn := fnan(x)
	if fm == 0 || (fn|fi) != 0 {
		return x
	}
	x &= uint64(0x7fffffffffffffff) //clear sign
	y := fmod(x)
	if fs != 0 {
		if fsub(x, y) != uint64(0) {
			y = fadd(y, uint64(0x3ff0000000000000)) // y+1.
		}
	}
	return y | fs
}
func fmod(x uint64) uint64 { // x>0.
	if x == 0 {
		return x
	}
	if x < uint64(0x3ff0000000000000) { // x<1.0
		return 0
	}
	e := (x>>uint64(52))&2047 - 1023
	if e < 52 {
		x &^= uint64(1)<<(52-e) - 1
	}
	return x
}
func fcps(x, y uint64) uint64 { return x&^uint64(0x8000000000000000) | y&uint64(0x8000000000000000) }
func fmax(x, y uint64) uint64 {
	if fnan(x) != 0 {
		return y
	}
	if fnan(y) != 0 {
		return x
	}
	xs := fsgn(x)
	ys := fsgn(y)
	if xs != ys {
		if xs != 0 {
			return y
		}
		return x
	}
	if x < y {
		return y
	}
	return x
}
func fmin(x, y uint64) uint64 {
	if fnan(x) != 0 {
		return y
	}
	if fnan(y) != 0 {
		return x
	}
	xs := fsgn(x)
	ys := fsgn(y)
	if xs != ys {
		if xs != 0 {
			return x
		}
		return y
	}
	if x < y {
		return x
	}
	return y
}
func fsqr(x uint64) uint64 { // see go.dev/src/math/sqrt.go
	if x&uint64(0x7fffffffffffffff) == 0 || fnan(x) != 0 || x == uint64(0x7ff0000000000000) {
		return x
	}
	if fsgn(x) != 0 {
		return nau
	}
	e := int32(x>>52) & 2047
	if e == 0 {
		for x&0x10000000000000 == 0 {
			x <<= uint64(1)
			e--
		}
		e++
	}
	e -= 1023
	x &^= 0x7ff0000000000000
	x |= 0x10000000000000
	if e&1 == 1 {
		x <<= uint64(1)
	}
	e >>= 1
	x <<= uint64(1)
	var q, s uint64
	r := uint64(0x20000000000000)
	for r != 0 {
		t := s + r
		if t <= x {
			s = t + r
			x -= t
			q += r
		}
		x <<= uint64(1)
		r >>= uint64(1)
	}
	if x != 0 {
		q += q & uint64(1)
	}
	return q>>uint64(1) + uint64(e+1022)<<uint64(52)
}
