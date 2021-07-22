package main

import (
	. "github.com/ktye/wg/module"
)

const pi float64 = 3.141592653589793
const maxfloat float64 = 1.797693134862315708145274237317043567981e+308

func hypot(p, q float64) float64 {
	//todo
	//switch {
	//case IsInf(p, 0) || IsInf(q, 0):
	//	return Inf(1)
	//case IsNaN(p) || IsNaN(q):
	//	return NaN()
	//}
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

func cosin(deg float64) (c, s float64) {
	if deg == 0 {
		c = 1.0
	} else if deg == 90 {
		s = 1.0
	} else if deg == 180 {
		c = -1.0
	} else if deg == 270 {
		s = -1.0
	} else {
		c, s = cosin_(deg * 0.017453292519943295)
	}
	return c, s
}
func cosin_(x float64) (c, s float64) {
	ss, cs := false, false
	if x < 0 {
		x = -x
		ss = true
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
		ss, cs = !ss, !cs
	}
	if j > 1 {
		cs = !cs
	}
	zz := z * z
	c = 1.0 - 0.5*zz + zz*zz*((((((-1.13585365213876817300e-11*zz)+2.08757008419747316778e-9)*zz+-2.75573141792967388112e-7)*zz+2.48015872888517045348e-5)*zz+-1.38888888888730564116e-3)*zz+4.16666666666665929218e-2)
	s = z + z*zz*((((((1.58962301576546568060e-10*zz)+-2.50507477628578072866e-8)*zz+2.75573136213857245213e-6)*zz+-1.98412698295895385996e-4)*zz+8.33333333332211858878e-3)*zz+-1.66666666666666307295e-1)
	if j == 1 || j == 2 {
		x = c
		c = s
		s = x
	}
	if cs {
		c = -c
	}
	if ss {
		s = -s
	}
	return c, s
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
	if x == 0 {
		return x
	}
	if x > 0 {
		return satan(x)
	}
	return -satan(-x)
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
func signbit(x float64) int32 { return int32(I64reinterpret_f64(x) >> 63) }
func exp(x float64) float64 {
	if isnan(x) {
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
func ldexp(frac float64, exp int64) float64 {
	if frac == 0 || frac > maxfloat || frac < -maxfloat || isnan(frac) {
		return frac
	}
	var e int64
	frac, e = normalize(frac)
	exp += e
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
func frexp(f float64) (frac float64, exp int64) {
	if f == 0.0 {
		return f, 0
	}
	if f < -maxfloat || f > maxfloat || isnan(f) {
		return f, 0
	}
	f, exp = normalize(f)
	x := I64reinterpret_f64(f)
	exp += int64((x>>52)&2047) - 1022
	x &^= 9218868437227405312
	x |= 4602678819172646912
	return F64reinterpret_i64(x), exp
}

func normalize(x float64) (y float64, exp int64) {
	if F64abs(x) < 2.2250738585072014e-308 {
		return x * 4.503599627370496e+15, int64(-52)
	}
	return x, 0
}

func log(x float64) float64 {
	if isnan(x) || x > maxfloat {
		return x
	}
	if x < 0 {
		return nan
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
func lgn(x, y float64) float64 {

	return log(y) * 1.0 / log(x)
}

func modf(f float64) (i float64, frac float64) {
	if f < 1.0 {
		if f < 0.0 {
			i, f = modf(-f)
			return -i, -f
		}
		if f == 0.0 {
			return f, f
		}
		return 0, f
	}
	x := I64reinterpret_f64(f)
	e := (x>>52)&2047 - 1023
	if e < 52 {
		x &^= 1<<(52-e) - 1
	}
	i = F64reinterpret_i64(x)
	frac = f - i
	return i, frac
}
func pow(x, y float64) float64 {
	if y == 0.0 || x == 1.0 {
		return 1.0
	}
	if y == 1.0 {
		return x
	}
	if isnan(x) || isnan(y) || y > maxfloat || y < -maxfloat { // simplified
		return nan
	}
	if x == 0 {
		if y < 0 {
			if isOddInt(y) {
				return F64copysign(inf, x)
			} else {
				return inf
			}
		} else {
			if isOddInt(y) {
				return x
			} else {
				return 0.0
			}
		}
	}
	if y == 0.5 {
		return F64sqrt(x)
	}
	if y == -0.5 {
		return 1.0 / F64sqrt(x)
	}

	yi, yf := modf(F64abs(y))
	if yf != 0.0 && x < 0.0 {
		return nan
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
func isOddInt(x float64) bool {
	xi, xf := modf(x)
	return xf == 0 && int64(xi)&1 == 1
}

/*
func exp(x float64) float64 {
	const (
		Ln2Hi = 6.93147180369123816490e-01
		Ln2Lo = 1.90821492927058770002e-10
		Log2e = 1.44269504088896338700e+00

		Overflow  = 7.09782712893383973096e+02
		Underflow = -7.45133219101941108420e+02
		NearZero  = 1.0 / (1 << 28) // 2**-28
	)
	switch {
	case IsNaN(x) || IsInf(x, 1):
		return x
	case IsInf(x, -1):
		return 0
	case x > Overflow:
		return Inf(1)
	case x < Underflow:
		return 0
	case -NearZero < x && x < NearZero:
		return 1 + x
	}
	var k int
	switch {
	case x < 0:
		k = int(Log2e*x - 0.5)
	case x > 0:
		k = int(Log2e*x + 0.5)
	}
	hi := x - float64(k)*Ln2Hi
	lo := float64(k) * Ln2Lo
	return expmulti(hi, lo, k)
}
func expmulti(hi, lo float64, k int) float64 {
	const (
		P1 = 1.66666666666666657415e-01
		P2 = -2.77777777770155933842e-03
		P3 = 6.61375632143793436117e-05
		P4 = -1.65339022054652515390e-06
		P5 = 4.13813679705723846039e-08
	)

	r := hi - lo
	t := r * r
	c := r - t*(P1+t*(P2+t*(P3+t*(P4+t*P5))))
	y := 1 - ((lo - (r*c)/(2-c)) - hi)
	return Ldexp(y, k)
}
func ldexp(frac float64, exp int) float64 {
	// special cases
	switch {
	case frac == 0:
		return frac // correctly return -0
	case IsInf(frac, 0) || IsNaN(frac):
		return frac
	}
	frac, e := normalize(frac)
	exp += e
	x := Float64bits(frac)
	exp += int(x>>shift)&mask - bias
	if exp < -1075 {
		return Copysign(0, frac) // underflow
	}
	if exp > 1023 { // overflow
		if frac < 0 {
			return Inf(-1)
		}
		return Inf(1)
	}
	var m float64 = 1
	if exp < -1022 { // denormal
		exp += 53
		m = 1.0 / (1 << 53) // 2**-53
	}
	x &^= mask << shift
	x |= uint64(exp+bias) << shift
	return m * Float64frombits(x)
}

*/

/*
package sincos (from go pkg math)

import (
	"math"
	"unsafe"
)

func Sqrt(f float64) float64    { return math.Sqrt(f) }
func IsNaN(f float64) (is bool) { return f != f }
func NaN() float64              { return Float64frombits(uvnan) }
func IsInf(f float64, sign int) bool {
	return sign >= 0 && f > MaxFloat64 || sign <= 0 && f < -MaxFloat64
}
func Float64frombits(b uint64) float64 { return *(*float64)(unsafe.Pointer(&b)) }
func Float64bits(f float64) uint64     { return *(*uint64)(unsafe.Pointer(&f)) }
func Abs(x float64) float64            { return Float64frombits(Float64bits(x) &^ (1 << 63)) }
func Inf(sign int) float64 {
	var v uint64
	if sign >= 0 {
		v = uvinf
	} else {
		v = uvneginf
	}
	return Float64frombits(v)
}
func Signbit(x float64) bool { return Float64bits(x)&(1<<63) != 0 }
func Copysign(x, y float64) float64 {
	const sign = 1 << 63
	return Float64frombits(Float64bits(x)&^sign | Float64bits(y)&sign)
}

var _sin = [...]float64{
	1.58962301576546568060e-10, // 0x3de5d8fd1fd19ccd
	-2.50507477628578072866e-8, // 0xbe5ae5e5a9291f5d
	2.75573136213857245213e-6,  // 0x3ec71de3567d48a1
	-1.98412698295895385996e-4, // 0xbf2a01a019bfdf03
	8.33333333332211858878e-3,  // 0x3f8111111110f7d0
	-1.66666666666666307295e-1, // 0xbfc5555555555548
}
var _cos = [...]float64{
	-1.13585365213876817300e-11, // 0xbda8fa49a0861a9b
	2.08757008419747316778e-9,   // 0x3e21ee9d7b4e3f05
	-2.75573141792967388112e-7,  // 0xbe927e4f7eac4bc6
	2.48015872888517045348e-5,   // 0x3efa01a019c844f5
	-1.38888888888730564116e-3,  // 0xbf56c16c16c14f91
	4.16666666666665929218e-2,   // 0x3fa555555555554b
}

const (
	uvinf      = 0x7FF0000000000000
	uvneginf   = 0xFFF0000000000000
	uvnan      = 0x7FF8000000000001
	MaxFloat64 = 1.797693134862315708145274237317043567981e+308
	Pi         = 3.14159265358979323846264338327950288419716939937510582097494459
	PI4A       = 7.85398125648498535156e-1  // 0x3fe921fb40000000, Pi/4 split into three parts
	PI4B       = 3.77489470793079817668e-8  // 0x3e64442d00000000,
	PI4C       = 2.69515142907905952645e-15 // 0x3ce8469898cc5170,
)

func sincos(x float64) (sin, cos float64) {
	switch {
	case x == 0:
		return x, 1 // ±0.0
	case IsNaN(x) || IsInf(x, 0):
		return NaN(), NaN()
	}
	sinSign, cosSign := false, false
	if x < 0 {
		x = -x
		sinSign = true
	}
	var j uint64
	var y, z float64
	j = uint64(x * (4 / Pi))
	y = float64(j)
	if j&1 == 1 { // map zeros to origin
		j++
		y++
	}
	j &= 7
	z = ((x - y*PI4A) - y*PI4B) - y*PI4C
	if j > 3 {
		j -= 4
		sinSign, cosSign = !sinSign, !cosSign
	}
	if j > 1 {
		cosSign = !cosSign
	}
	zz := z * z
	cos = 1.0 - 0.5*zz + zz*zz*((((((_cos[0]*zz)+_cos[1])*zz+_cos[2])*zz+_cos[3])*zz+_cos[4])*zz+_cos[5])
	sin = z + z*zz*((((((_sin[0]*zz)+_sin[1])*zz+_sin[2])*zz+_sin[3])*zz+_sin[4])*zz+_sin[5])
	if j == 1 || j == 2 {
		sin, cos = cos, sin
	}
	if cosSign {
		cos = -cos
	}
	if sinSign {
		sin = -sin
	}
	return
}
func hypot(p, q float64) float64 {
	switch {
	case IsInf(p, 0) || IsInf(q, 0):
		return Inf(1)
	case IsNaN(p) || IsNaN(q):
		return NaN()
	}
	p, q = Abs(p), Abs(q)
	if p < q {
		p, q = q, p
	}
	if p == 0 {
		return 0
	}
	q = q / p
	return p * Sqrt(1+q*q)
}
func atan2(y, x float64) float64 {
	switch {
	case IsNaN(y) || IsNaN(x):
		return NaN()
	case y == 0:
		if x >= 0 && !Signbit(x) {
			return Copysign(0, y)
		}
		return Copysign(Pi, y)
	case x == 0:
		return Copysign(Pi/2, y)
	case IsInf(x, 0):
		if IsInf(x, 1) {
			switch {
			case IsInf(y, 0):
				return Copysign(Pi/4, y)
			default:
				return Copysign(0, y)
			}
		}
		switch {
		case IsInf(y, 0):
			return Copysign(3*Pi/4, y)
		default:
			return Copysign(Pi, y)
		}
	case IsInf(y, 0):
		return Copysign(Pi/2, y)
	}
	q := atan(y / x)
	if x < 0 {
		if q <= 0 {
			return q + Pi
		}
		return q - Pi
	}
	return q
}
func atan(x float64) float64 {
	if x == 0 {
		return x
	}
	if x > 0 {
		return satan(x)
	}
	return -satan(-x)
}
func satan(x float64) float64 {
	const (
		Morebits = 6.123233995736765886130e-17 // pi/2 = PIO2 + Morebits
		Tan3pio8 = 2.41421356237309504880      // tan(3*pi/8)
	)
	if x <= 0.66 {
		return xatan(x)
	}
	if x > Tan3pio8 {
		return Pi/2 - xatan(1/x) + Morebits
	}
	return Pi/4 + xatan((x-1)/(x+1)) + 0.5*Morebits
}
func xatan(x float64) float64 {
	const (
		P0 = -8.750608600031904122785e-01
		P1 = -1.615753718733365076637e+01
		P2 = -7.500855792314704667340e+01
		P3 = -1.228866684490136173410e+02
		P4 = -6.485021904942025371773e+01
		Q0 = +2.485846490142306297962e+01
		Q1 = +1.650270098316988542046e+02
		Q2 = +4.328810604912902668951e+02
		Q3 = +4.853903996359136964868e+02
		Q4 = +1.945506571482613964425e+02
	)
	z := x * x
	z = z * ((((P0*z+P1)*z+P2)*z+P3)*z + P4) / (((((z+Q0)*z+Q1)*z+Q2)*z+Q3)*z + Q4)
	z = x*z + x
	return z
}
*/
