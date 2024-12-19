//go:build simd5

package main

//gcc 32byte vector extensions(avx2..)

import (
        . "github.com/ktye/wg/module"
)

const b0 = 6 //smallest bucket
const bs = 31 //lower bits set
const vl = 32 //vector length

func ev(ep int32) int32 { return (31 + ep) & -32 }
func seq(n int32) K {
        i := Iota5()
        n = maxi(n, 0)
        r := mk(It, n)
        rp := int32(r)
        e := rp + ev(4*n)
        for rp < e {
                I5store(rp, i)
                i = i.Add(I5splat(8))
                rp += 32
        }
        return r
}

