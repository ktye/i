// +build:ignore
package main

// why? be kos

/*
func fneg(x float64) float64 { return 0 }
func fsqt(x float64) float64 { return 0 }
func fadd(x, y float64) float64 { return 0 }
func fsub(x, y float64) float64 { return 0 }
func fmul(x, y float64) float64 { return 0 }
func fdiv(x, y float64) float64 { return 0 }

func fsgn(x uint64) int32 { return int32(x>>32)&(1<<32) }
func fmnt(x uint64) uint64 { return x & 4503599627370495 }
func fexp(x uint64) int32 { return int32(x>>52) & 2047 }
func fnan(x uint64) int32 { return I32B(fexp(x) == 2047) }
*/

/*
mantbits64 uint = 52
expbits64  uint = 11
bias64          = -1<<(expbits64-1) + 1
nan64 uint64 = (1<<expbits64-1)<<mantbits64 + 1<<(mantbits64-1) // quiet NaN, 0 payload
inf64 uint64 = (1<<expbits64 - 1) << mantbits64
neg64 uint64 = 1 << (expbits64 + mantbits64)
*/


// https://go.dev/src/runtime/softfloat64.go
