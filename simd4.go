//go:build simd4

package main // wasm simd128 (16bytes)

import (
	. "github.com/ktye/wg/module"
)

const b0 = 5 //smallest bucket
const bs = 15 //lower bits set
const vl = 16 //vector length 

func ev(ep int32) int32 { return (15 + ep) & -16 }
func seq(n int32) K {
	i := Iota4()
	n = maxi(n, 0)
	r := mk(It, n)
	rp := int32(r)
	e := rp + ev(4*n)
	for rp < e {
		I4store(rp, i)
		i = i.Add(I4splat(4))
		rp += 16
	}
	return r
}
