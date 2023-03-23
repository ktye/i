package main

// why? be kos

// from go.dev/src/runtime/softfloat64.go

func fsgn(x uint64) uint64 { return x & 9223372036854775808 }
func fmnt(x uint64) uint64 { return x & 4503599627370495 }
func fexp(x uint64) int32  { return int32(x>>52) & 2047 }
func fnan(x uint64) int32  { return I32B(fexp(x) == 2047 && fmnt(x) != 0) }
func finf(x uint64) int32  { return I32B(fexp(x) == 2047 && fmnt(x) == 0) }

func fpak(s, m uint64, e int32, t uint64) float64 {
	// todo
	return F64reinterpret_i64(s)
}

func fadd(xf, yf float64) float64 {
	var ti int32
	var tu uint64
	x := I64reinterpret_f64(xf)
	y := I64reinterpret_f64(yf)
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
		return na
	}
	if fi != 0 && gi != 0 && fs != gs {
		return na
	}
	if fi != 0 {
		return xf
	}
	if gi != 0 {
		return yf
	}
	if fm == 0 && gm == 0 && fs != 0 && gs != 0 {
		return xf
	}
	if fm == 0 {
		if gm == 0 {
			y ^= fs
		}
		return F64reinterpret_i64(y)
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
	return fpak(fs, fm, fe-2, t)
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
