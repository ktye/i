//go:build small

package main

import (
	. "github.com/ktye/wg/module"
)

const SMALL = true

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

	//Functions(367, sum, rd0, prd, rd0, min, max)
	//Functions(374, sums, rd0, prds, rd0, rd0)
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
