//go:build !(simd4||simd5)

package main

import (
	. "github.com/ktye/wg/module"
)

const b0 = 5 //smallest bucket
const bs = 15 //lower bits set
const vl = 16 //vector length

func seq(n int32) K {
        n = maxi(n, 0)
        r := mk(It, n)
        for n > 0 {
                n--
                SetI32(int32(r)+4*n, n)
        }
        return r
}

