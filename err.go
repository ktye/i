package k

import . "github.com/ktye/wg/module"

const (
	Err int32 = iota
	Type
	Value
	Length
	Rank
	Parse
	Grow
	Unref
	Nyi
)

func trap(x int32) { SetI32(0, x); panic(x) }
