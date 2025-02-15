//go:build simd5

package main

//gcc 32byte vector extensions(avx2..)

import (
	. "github.com/ktye/wg/module"
)

func init() {
	Simd(5)
}

const b0 = 6 //smallest bucket
const bs = 31 //lower bits set
const vl = 32 //vector length

func ev(ep int32) int32 { return (31 + ep) & -32 }
func eu(ep int32) int32 { return ep & -32 }
func sumz(x, e, r int32) {
	a := VFsplat(0.0)
	u := eu(e)
	for x < u {
		a = a.Add(VFload(x))
		x += vl
	}
	re := a.HsumEven()
	im := a.HsumOdd()
	for x < e {
		re += F64(x)
		im += F64(x+8)
		x += 16
	}
	SetF64(r, re)
	SetF64(r+8, im)
}
