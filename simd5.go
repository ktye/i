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
