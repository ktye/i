//go:build small

package main

import (
	. "github.com/ktye/wg/module"
)

const SMALL = true

func init1() {
	Functions(220, negi, negf, negz)
	Functions(227, absi, absf, absz)
	Functions(234, addi, addf, addz)
	Functions(245, subi, subf, subz)
	Functions(256, muli, mulf, mulz)
	Functions(267, divi, divf, divz, modi, nyi, nyi)
	Functions(278, mini, minf, minz)
	Functions(289, maxi, maxf, maxz)
	Functions(300, nyi, sqrf, nyi)

	Functions(308, lti, ltf, ltz)
	Functions(323, gti, gtf, gtz)
	Functions(338, eqi, eqf, eqz)

	Functions(353, guC, guC, guI, guI, guF, guZ, guL, gdC, gdC, gdI, gdI, gdF, gdZ, gdL)

}

/* use only this syntax:
if SMALL {
	...
}
if SMALL == false {
	...
}
*/

func ov0(f, x K) (r K) { // f/0#x
	r = missing(tp(x) - 16)
	if tp(f) == 0 { // 3174190 is `".oO"
		if int32(f) == 13 {
			return x
		}
		dx(x)
		//println("f-2", int32(f-2), tp(r))
		return Cal(Val(sc(Ku(3174190))), l2(Ki(int32(f-2)), r))
	}
	dx(f)
	dx(x)
	return r
}

func min(x K, yp int32, t T, n int32) K  { return Rdc(6, rx(K(uint32(x))|K(t)<<59)) }
func max(x K, yp int32, t T, n int32) K  { return Rdc(7, rx(K(uint32(x))|K(t)<<59)) }
func epx(f int32, x, y K, n int32) (r K) { return r }
func epc(f int32, x, y K, n int32) (r K) { return r }
func fndXs(x, y K, t T, yn int32) (r K)  { return r }
