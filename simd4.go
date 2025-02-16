//go:build simd4

package main

// 16byte simd, e.g. wasm simd128

import (
	. "github.com/ktye/wg/module"
)

func init() {
	Simd(4)
}

const b0 = 5  //smallest bucket
const bs = 15 //lower bits set
const vl = 16 //vector length

func ev(ep int32) int32 { return (15 + ep) & -16 }
func eu(ep int32) int32 { return ep & -16 }
func sumz(x, e, r int32) {
	a := VFsplat(0.0)
	for x < e {
		a = a.Add(VFload(x))
		x += vl
	}
	VFstore(r, a)
}
