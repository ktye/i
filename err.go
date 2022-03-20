package main

import (
	. "github.com/ktye/wg/module"
)

const (
	Err int32 = iota
	Type
	Value
	Index
	Length
	Rank
	Parse
	Stack
	Grow
	Unref
	Io
	Nyi
)

func trap(x int32) K {
	src := src()
	//Printf("src %d %d, nn=%d srcp:%d\n", int32(src), int32(src>>32), nn(src), srcp)
	//if srcp < nn(src) {
	if srcp == 0 {
		write(Ku(2608)) // 0\n
	} else {
		a := maxi(srcp-30, 0)
		for i := a; i < srcp; i++ {
			if I8(int32(src)+i) == 10 {
				a = 1 + i
			}
		}
		b := mini(int32(src)+nn(src), srcp+30)
		for i := srcp; i < b; i++ {
			if I8(int32(src)+i) == 10 {
				b = i
				break
			}
		}
		Write(0, 0, int32(src)+a, b-a)
		if srcp > a {
			write(Cat(Kc(10), ntake(srcp-a-1, Kc(32))))
		}
	}
	write(Ku(2654)) // ^\n
	panic(x)
	return 0
}
func Srcp() int32 { return srcp }
