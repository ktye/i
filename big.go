//go:build !small

package main

import (
	. "github.com/ktye/wg/module"
)

func init1() {
	Functions(220, negi, negf, negz, negC, negI, negF, negZ)
	Functions(227, absi, absf, nyi, absC, absI, absF, absZ)
	Functions(234, addi, addf, addz, addcC, addiI, addfF, addzZ, addC, addI, addF, addF)
	Functions(245, subi, subf, nyi, subcC, subiI, subfF, subzZ, subC, subI, subF, subF)
	Functions(256, muli, mulf, mulz, mulcC, muliI, mulfF, mulzZ, mulC, mulI, mulF, mulZ)
	Functions(267, divi, divf, divz, nyi, nyi, divfF, divzZ, nyi, nyi, divF, divZ)
	Functions(278, mini, minf, minz, mincC, miniI, minfF, minzZ, minC, minI, minF, minZ)
	Functions(289, maxi, maxf, maxz, maxcC, maxiI, maxfF, maxzZ, maxC, maxI, maxF, maxZ)
	Functions(300, nyi, sqrf, nyi, nyi, nyi, sqrF, nyi)

	Functions(308, lti, ltf, ltz, ltcC, ltiI, ltfF, ltzZ, ltCc, ltIi, ltFf, ltZz, ltC, ltI, ltF, ltZ)
	Functions(323, gti, gtf, gtz, gtcC, gtiI, gtfF, gtzZ, gtCc, gtIi, gtFf, gtZz, gtC, gtI, gtF, gtZ)
	Functions(338, eqi, eqf, eqz, eqcC, eqiI, eqfF, eqzZ, eqCc, eqIi, eqFf, eqZz, eqC, eqI, eqF, eqZ)

	Functions(353, guC, guC, guI, guI, guF, guZ, guL, gdC, gdC, gdI, gdI, gdF, gdZ, gdL)

	Functions(367, sum, rd0, prd, rd0, min, max)
	Functions(374, sums, rd0, prds, rd0, rd0)
	// don't delete: kom:FTAB 381)

}

const SMALL = false

func ov0(f, x K) K {
	dx(f)
	dx(x)
	return missing(tp(x))
}

func epx(f int32, x, y K, n int32) (r K) { // ( +-*%&| )':
	xt := tp(x)
	xp := int32(x)
	r = mk(xt, n)
	rp := int32(r)
	f = 212 + 11*f
	yp := int32(y)
	if xt == It {
		SetI32(rp, Func[f].(f2i)(I32(xp), yp))
		for i := int32(1); i < n; i++ {
			xp += 4
			rp += 4
			SetI32(rp, Func[f].(f2i)(I32(xp), I32(xp-4)))
		}
	} else {
		f++
		SetF64(rp, Func[f].(f2f)(F64(xp), F64(yp)))
		for i := int32(1); i < n; i++ {
			xp += 8
			rp += 8
			SetF64(rp, Func[f].(f2f)(F64(xp), F64(xp-8)))
		}
	}
	dx(x)
	dx(y)
	return r
}
func epc(f int32, x, y K, n int32) (r K) { // ( <>= )':
	xt := tp(x)
	xp := int32(x)
	s := sz(xt)
	r = mk(It, n)
	rp := int32(r)
	f = 188 + 15*f
	switch s >> 2 {
	case 0:
		SetI32(rp, Func[f].(f2i)(I8(xp), int32(y)))
		for i := int32(1); i < n; i++ {
			xp++
			rp += 4
			SetI32(rp, Func[f].(f2i)(I8(xp), I8(xp-1)))
		}
	case 1:
		SetI32(rp, Func[f].(f2i)(I32(xp), int32(y)))
		for i := int32(1); i < n; i++ {
			xp += 4
			rp += 4
			SetI32(rp, Func[f].(f2i)(I32(xp), I32(xp-4)))
		}
	default:
		f++
		SetI32(rp, Func[f].(f2c)(F64(xp), F64(int32(y))))
		for i := int32(1); i < n; i++ {
			xp += 8
			rp += 4
			SetI32(rp, Func[f].(f2c)(F64(xp), F64(xp-8)))
		}
	}
	dx(x)
	dx(y)
	return r
}

func fndXs(x, y K, t T, yn int32) (r K) {
	xn := nn(x)
	a := int32(min(0, int32(x), t, xn))
	b := 1 + (int32(max(0, int32(x), t, xn))-a)>>(3*I32B(t == St))
	if b > 256 && b > yn {
		return 0
	}
	if t == St {
		x, y = Div(Flr(x), Ki(8)), Div(Flr(y), Ki(8))
		a >>= 3
	}
	r = ntake(b, Ki(nai))
	rp := int32(r) - 4*a
	x0 := int32(x)
	xp := ep(x)
	if t == Ct {
		for xp > x0 {
			xp--
			SetI32(rp+4*I8(xp), xp-x0)
		}
	} else {
		for xp > x0 {
			xp -= 4
			SetI32(rp+4*I32(xp), (xp-x0)>>2)
		}
	}
	dx(x)
	return Atx(r, Add(Ki(-a), y))
}
