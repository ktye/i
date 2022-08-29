package main

import (
	. "github.com/ktye/wg/module"
)

const ( // iota is not supported by wg
	Err    = int32(0)
	Type   = int32(1)
	Value  = int32(2)
	Index  = int32(3)
	Length = int32(4)
	Rank   = int32(5)
	Parse  = int32(6)
	Stack  = int32(7)
	Grow   = int32(8)
	Unref  = int32(9)
	Io     = int32(10)
	Nyi    = int32(11)
)

func trap(x int32) K {
	s := src()
	//kdb:Trap(p,x,srcp,int64(s))
	//Printf("src %d %d, nn=%d srcp:%d\n", int32(src), int32(src>>32), nn(src), srcp)
	//if srcp < nn(src) {
	if srcp == 0 {
		write(Ku(2608)) // 0\n
	} else {
		a := maxi(srcp-30, 0)
		for i := a; i < srcp; i++ {
			if I8(int32(s)+i) == 10 {
				a = 1 + i
			}
		}
		b := mini(nn(s), srcp+30)
		for i := srcp; i < b; i++ {
			if I8(int32(s)+i) == 10 {
				b = i
				break
			}
		}
		Write(0, 0, int32(s)+a, b-a)
		if srcp > a {
			write(Cat(Kc(10), ntake(srcp-a-1, Kc(32))))
		}
	}
	write(Ku(2654)) // ^\n
	panic(x)
	return 0
}
func Srcp() int32 { return srcp }
