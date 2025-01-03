package main

import (
	. "github.com/ktye/wg/module"
)

func rd0(yp int32, t T, n int32) K { return 0 }
func min(yp int32, t T, e int32) K { // &/x
	var xp int32
	switch t - 18 {
	case 0: // Ct
		xp = 127
		for yp < e {
			xp = mini(xp, I8(yp))
			yp++
		}
		return Kc(xp)
	case 1: // It
		return Ki(minis(2147483647, yp, e))
	case 2: // St
		return Ks(minis((nn(K(I64(8))) << 3) - 8, yp, e))
	case 3: // Ft
		return Kf(minfs(yp, e))
	default:
		return 0
	}
}
func max(yp int32, t T, e int32) K { // |/x
	var xp int32
	switch t - 18 {
	case 0: // Ct
		xp = -128
		for yp < e {
			xp = maxi(xp, I8(yp))
			yp++
		}
		return Kc(xp)
	case 1: // It
		return Ki(maxis(nai, yp, e))
	case 2: // St
		return Ks(maxis(0, yp, e))
	case 3: // Ft
		return Kf(maxfs(yp, e))
	default:
		return 0
	}
}
func sum(yp int32, t T, e int32) K { // +/x
	xp := int32(0)
	switch t - 18 {
	case 0: // Ct
		for yp < e {
			xp += I8(yp)
			yp++
		}
		return Kc(xp)
	case 1: // It
		return Ki(xp + sumi(yp, e))
	case 2: // St
		return 0
	case 3: // Ft
		f := 0.0
		return Kf(f + sumf(yp, e, 8))
	case 4: // Zt
		re := 0.0
		im := 0.0
		return Kz(re+sumf(yp, e, 16), im+sumf(yp+8, e, 16))
	default:
		return 0
	}
}
func prd(yp int32, t T, e int32) K { // */x
	xp := int32(1)
	switch t - 18 {
	case 0: // Ct
		for yp < e {
			xp *= I8(yp)
			yp++
		}
		return Kc(xp)
	case 1: // It
		for yp < e {
			xp *= I32(yp)
			yp += 4
		}
		return Ki(xp)
	case 2: // St
		return 0
	case 3: // Ft
		f := 1.0
		for yp < e {
			f *= F64(yp)
			yp += 8
		}
		return Kf(f)
	default:
		return 0
	}
}
