package main

import (
	. "github.com/ktye/wg/module"
)

const pi float64 = 3.141592653589793

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
