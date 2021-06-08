package main

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

func trap(x int32) {
	if src != 0 {
		write(cat1(src, Kc(10))) // todo clip
		if srcp != 0 {
			write(ntake(srcp-1, Kc(32)))
			write(Ku(2654)) // ^\n
		}
	}
	panic(x)
}
func Srcp() int32 { return srcp }
