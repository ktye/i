package main

import . "github.com/ktye/wg/module"

const (
	Err int32 = iota
	Type
	Value
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
	if src != 0 {

		a := maxi(srcp-30, 0)
		for i := a; i < srcp; i++ {
			if I8(int32(src)+i) == 10 {
				a = 1 + i
			}
		}
		b := mini(nn(src), srcp+30)
		for i := srcp; i < b; i++ {
			if I8(int32(src)+i) == 10 {
				b = i
				break
			}
		}

		src = cat1(ntake(b-a, ndrop(a, src)), Kc(10))
		write(src)

		srcp -= a
		if srcp > 0 {
			write(ntake(srcp-1, Kc(32)))
		}
		write(Ku(2654)) // ^\n
		srcp += a
	}
	panic(x)
	return 0
}
func Srcp() int32 { return srcp }
